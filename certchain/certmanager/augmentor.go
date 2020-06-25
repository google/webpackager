// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certmanager

import (
	"log"
	"sync"
	"time"

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"github.com/google/webpackager/internal/chanutil"
	"github.com/google/webpackager/internal/timeutil"
)

// Augmentor combines RawChainSource and OCSPRespSource to serve as a Producer.
type Augmentor struct {
	rawChain *certchain.RawChain
	rcSource RawChainSource
	ocspResp *certchain.OCSPResponse
	orSource OCSPRespSource
	outMu    sync.RWMutex
	out      chan *certchain.AugmentedChain
	killer   *chanutil.Killer
}

var _ Producer = (*Augmentor)(nil)

// RawChainSource provides a RawChain. It is designed to be called repeatedly:
// the Fetch method does not just return the latest certificate chain, but
// also instructs when it should be called again to receive the next update.
type RawChainSource interface {
	// Fetch returns a new RawChain to replace chain with a futureevent.Event
	// to notify when the RawChainSource expects the next call of Fetch.
	//
	// chain can be nil. Fetch returns a valid RawChain, which the caller can
	// use as the initial one.
	//
	// now is a function which returns the current time. Fetch calls it after
	// retrieving and parsing the certificate chain and examines the validity
	// as of the returned time. time.Now will do in most of the time, but the
	// caller may want to use a different function in unit testing or to deal
	// with clock skews.
	//
	// Fetch may return chain as newChain if it is valid and still up-to-date.
	//
	// nextRun is always non-nil and valid, even in the error case.
	Fetch(chain *certchain.RawChain, now func() time.Time) (newChain *certchain.RawChain, nextRun futureevent.Event, err error)
}

// OCSPRespSource provides an OCSPResponse. It is designed to be called
// repeatedly: the Fetch method does not just return the OCSP response, but
// also instructs when it should be called again to receive the next update.
type OCSPRespSource interface {
	// Fetch retrieves an OCSP response for chain and returns the parsed OCSP
	// response with a futureevent.Event to notify when the OCSPRespSource
	// expects the next call of Fetch.
	//
	// now is a function which returns the current time. Fetch calls it after
	// retriving and parsing the OCSP response and examines the validity as
	// of the returned time. time.Now will do in most of the time, but the
	// caller may want to use a different function in unit testing or to deal
	// with clock skews.
	//
	// nextRun is always non-nil and valid, even in the error case.
	Fetch(chain *certchain.RawChain, now func() time.Time) (ocspResp *certchain.OCSPResponse, nextRun futureevent.Event, err error)
}

// NewAugmentor creates and initializes a new Augmentor. orSource can be nil,
// in which case DefaultOCSPClient is used.
//
// NewAugmentor does not start the production of AugmentedChains automatically.
// To start it, call Start.
func NewAugmentor(rcSource RawChainSource, orSource OCSPRespSource) *Augmentor {
	if orSource == nil {
		orSource = DefaultOCSPClient
	}
	return &Augmentor{rcSource: rcSource, orSource: orSource}
}

// Start spawns a goroutine to produce AugmentedChains continuously. It produces
// the first AugmentedChain before starting the goroutine and blocks until it is
// ready.
//
// The goroutine is not spawned in case of error.
func (a *Augmentor) Start() error {
	a.outMu.Lock()
	a.out = make(chan *certchain.AugmentedChain, 1)
	a.outMu.Unlock()

	a.killer = chanutil.NewKiller()

	rcNext, _, err := a.maintainRawChain()
	if err != nil {
		return err
	}
	orNext, err := a.maintainOCSPResp()
	if err != nil {
		return err
	}

	go a.daemon(rcNext, orNext)
	return nil
}

func (a *Augmentor) daemon(rcNext, orNext futureevent.Event) {
	for {
		var err error
		select {
		case <-rcNext.Chan():
			var updated bool
			rcNext, updated, err = a.maintainRawChain()
			if err != nil {
				log.Printf("cannot update the certificate: %v", err)
			}
			if updated {
				orNext.Cancel()
				orNext, err = a.maintainOCSPResp()
				if err != nil {
					log.Printf("cannot update the OCSP response: %v", err)
				}
			}
		case <-orNext.Chan():
			orNext, err = a.maintainOCSPResp()
			if err != nil {
				log.Printf("cannot update the OCSP response: %v", err)
			}
		case <-a.killer.C:
			rcNext.Cancel()
			orNext.Cancel()
			return
		}
	}
}

func (a *Augmentor) maintainRawChain() (nextRun futureevent.Event, updated bool, err error) {
	newChain, nextRun, err := a.rcSource.Fetch(a.rawChain, timeutil.Now)
	if err != nil {
		return nextRun, false, err
	}
	if a.rawChain == nil || a.rawChain.Digest != newChain.Digest {
		a.rawChain = newChain
		updated = true
	}
	return nextRun, updated, nil
}

func (a *Augmentor) maintainOCSPResp() (futureevent.Event, error) {
	ocspResp, nextRun, err := a.orSource.Fetch(a.rawChain, timeutil.Now)
	if err != nil {
		return nextRun, err
	}
	a.ocspResp = ocspResp

	a.outMu.RLock()
	a.out <- certchain.NewAugmentedChain(a.rawChain, a.ocspResp, nil)
	a.outMu.RUnlock()

	return nextRun, nil
}

// Stop kills the goroutine spawned by Start to stop producing AugmentedChains.
func (a *Augmentor) Stop() {
	a.outMu.Lock()
	defer a.outMu.Unlock()

	if a.out == nil {
		return
	}
	a.killer.Kill()
	close(a.out)
	a.out = nil
}

// Out returns the channel to receive produced AugmentedChains.
func (a *Augmentor) Out() <-chan *certchain.AugmentedChain {
	a.outMu.RLock()
	defer a.outMu.RUnlock()
	return a.out
}
