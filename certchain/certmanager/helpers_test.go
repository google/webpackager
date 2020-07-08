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

package certmanager_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

const (
	defaultTimeout = 3 * time.Second
	instantTimeout = time.Nanosecond
)

var (
	bigIntComparer = cmp.Comparer(func(x, y *big.Int) bool {
		return x.Cmp(y) == 0
	})

	certComparer = cmp.Options{bigIntComparer}
	ocspComparer = cmp.Options{bigIntComparer}
)

func createTempFile() string {
	f, err := ioutil.TempFile("", "certmanager_test_")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return f.Name()
}

func copyFile(srcFile, dstFile string) {
	src, err := os.Open(srcFile)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.OpenFile(dstFile, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	io.Copy(dst, src)
}

//------------------------
//  stubProducer

var _ certmanager.Producer = (*stubProducer)(nil)

type stubProducer struct {
	out chan *certchain.AugmentedChain
}

func newStubProducer() *stubProducer {
	return &stubProducer{make(chan *certchain.AugmentedChain, 1)}
}

func (p *stubProducer) Out() <-chan *certchain.AugmentedChain { return p.out }

func (p *stubProducer) Start() error { return nil }

func (p *stubProducer) Stop() { close(p.out) }

//------------------------
//  stubCache
// TODO(banaag): remove this and replace with in-memory cache.

var _ certmanager.Cache = (*stubCache)(nil)

type stubCache struct {
	avail    chan struct{}
	chainMap map[string]*certchain.AugmentedChain
}

func newStubCache() *stubCache {
	return &stubCache{
		make(chan struct{}, 1),
		make(map[string]*certchain.AugmentedChain),
	}
}

func (c *stubCache) Read(digest string) (*certchain.AugmentedChain, error) {
	if _, ok := c.chainMap[digest]; !ok {
		return nil, certmanager.ErrNotFound
	}
	return c.chainMap[digest], nil
}

func (c *stubCache) ReadLatest() (*certchain.AugmentedChain, error) {
	// TODO(banaag): Add support once we have in-memory cache.
	return nil, certmanager.ErrNotFound
}

func (c *stubCache) Write(ac *certchain.AugmentedChain) error {
	c.chainMap[ac.Digest] = ac
	c.avail <- struct{}{}
	return nil
}

func (c *stubCache) WaitForAvail(timeout time.Duration) waitResult {
	select {
	case _, ok := <-c.avail:
		if !ok {
			return waitCanceled
		}
		return waitSuccess
	case <-time.After(timeout):
		return waitTimeout
	}
}

//------------------------
//  stubRawChainSource

var _ certmanager.RawChainSource = (*stubRawChainSource)(nil)

type stubRawChainSource struct {
	data        []*certchain.RawChain
	NextRun     *futureevent.TriggerableEvent
	OnFetchDone chan time.Time
	mu          sync.RWMutex
}

func newStubRawChainSource(data ...*certchain.RawChain) *stubRawChainSource {
	return &stubRawChainSource{
		data:        data,
		OnFetchDone: make(chan time.Time, 1),
	}
}

func (s *stubRawChainSource) RestCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *stubRawChainSource) Fetch(chain *certchain.RawChain, now func() time.Time) (*certchain.RawChain, futureevent.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newChain := s.data[0]
	s.data = s.data[1:]
	s.NextRun = futureevent.NewTriggerableEvent()
	s.OnFetchDone <- now()

	return newChain, s.NextRun, nil
}

//------------------------
//  stubOCSPRespSource

var _ certmanager.OCSPRespSource = (*stubOCSPRespSource)(nil)

type stubOCSPRespSource struct {
	data        []*certchain.OCSPResponse
	NextRun     *futureevent.TriggerableEvent
	OnFetchDone chan time.Time
	mu          sync.RWMutex
}

func newStubOCSPRespSource(data ...*certchain.OCSPResponse) *stubOCSPRespSource {
	return &stubOCSPRespSource{
		data:        data,
		OnFetchDone: make(chan time.Time, 1),
	}
}

func (s *stubOCSPRespSource) RestCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *stubOCSPRespSource) Fetch(chain *certchain.RawChain, now func() time.Time) (*certchain.OCSPResponse, futureevent.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ocspResp := s.data[0]
	s.data = s.data[1:]
	s.NextRun = futureevent.NewTriggerableEvent()
	s.OnFetchDone <- now()

	return ocspResp, s.NextRun, nil
}

//------------------------
//  dummyFutureEvent

// dummyFutureEvent is a dummy Event, which never happens.
type dummyFutureEvent struct {
	Time time.Time
}

func newDummyFutureEvent(t time.Time) futureevent.Event {
	return &dummyFutureEvent{t}
}

func (*dummyFutureEvent) Chan() <-chan time.Time { return nil }

func (*dummyFutureEvent) Cancel() {}

//------------------------
//  waitResult

// waitResult represents the result of waitFor.
type waitResult int

const (
	waitSuccess waitResult = iota
	waitCanceled
	waitTimeout
)

func (w waitResult) String() string {
	switch w {
	case waitSuccess:
		return "waitSuccess"
	case waitCanceled:
		return "waitCanceled"
	case waitTimeout:
		return "waitTimeout"
	default:
		return fmt.Sprintf("waitResult(%d)", w)
	}
}

func waitFor(c <-chan time.Time, timeout time.Duration) waitResult {
	select {
	case _, ok := <-c:
		if !ok {
			return waitCanceled
		}
		return waitSuccess
	case <-time.After(timeout):
		return waitTimeout
	}
}

func waitForEvent(e futureevent.Event, timeout time.Duration) waitResult {
	return waitFor(e.Chan(), timeout)
}
