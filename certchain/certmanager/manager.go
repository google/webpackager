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

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/internal/chanutil"
)

// Config configures Manager.
type Config struct {
	// RawChainSource is used as part of a new Augmentor.
	// It may not be nil unless Producer is specified.
	RawChainSource RawChainSource

	// OCSPRespSource is used as part of a new Augmentor.
	// nil implies DefaultOCSPClient.
	OCSPRespSource OCSPRespSource

	// Producer allows specifying the Producer directly, especially using
	// a custom implementation. Producer takes precedence: RawChainSource
	// and OCSPRespSource will not be used if Producer is set non-nil.
	Producer Producer

	// Cache specifies where to cache the signed exchange certificates.
	// nil implies NullCache, i.e. no caching.
	Cache Cache
}

// Producer produces a new AugmentedChain repeatedly and sends it through
// a channel every time Producer completes the production.
type Producer interface {
	// Out returns the channel to receive produced AugmentedChains.
	Out() <-chan *certchain.AugmentedChain
	// Start starts producing new AugmentedChains.
	Start() error
	// Stop stops producing new AugmentedChains.
	Stop()
}

// Cache represents a storage to cache an AugmentedChain.
type Cache interface {
	Write(ac *certchain.AugmentedChain) error
}

// Manager keeps a signed exchange certificate up-to-date. See the package
// document for details.
type Manager struct {
	dataMu   sync.RWMutex
	data     *certchain.AugmentedChain
	producer Producer
	cache    Cache
	killer   *chanutil.Killer
}

// NewManager creates and initializes a new Manager.
func NewManager(c Config) *Manager {
	producer := c.Producer
	cache := c.Cache

	if producer == nil {
		producer = NewAugmentor(c.RawChainSource, c.OCSPRespSource)
	}
	if cache == nil {
		cache = NullCache
	}

	return &Manager{producer: producer, cache: cache}
}

// Start starts managing the certificate. It starts Producer and waits for
// the first AugmentedChain, then kicks in a goroutine to keep Cache updated
// with the received AugmentedChains. The execution is blocked until the first
// AugmentedChain comes in.
//
// To stop the management, call Stop.
func (m *Manager) Start() error {
	err := m.producer.Start()
	if err != nil {
		return err
	}

	m.data = <-m.producer.Out()
	go m.onReceive(m.data)

	m.killer = chanutil.NewKiller()
	go m.daemon()

	return nil
}

func (m *Manager) daemon() {
	for {
		select {
		case data := <-m.producer.Out():
			go m.onReceive(data)
			m.dataMu.Lock()
			m.data = data
			m.dataMu.Unlock()
		case <-m.killer.C:
			return
		}
	}
}

func (m *Manager) onReceive(ac *certchain.AugmentedChain) {
	if err := m.cache.Write(ac); err != nil {
		log.Printf("cannot cache the latest AugmentedChain: %v", err)
	}
}

// Stop stops the Manager m from managing the certificate: stops m.Producer
// and the cache updater goroutine spawned by Start.
//
// If m is writing an AugmentedChain to m.Cache when Stop is called, that
// cache write continues on background. Stop returns without waiting for its
// completion.
func (m *Manager) Stop() {
	m.killer.Kill()
	m.producer.Stop()
}

// GetAugmentedChain returns the AugmentedChain that m currently holds.
func (m *Manager) GetAugmentedChain() *certchain.AugmentedChain {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()
	return m.data
}
