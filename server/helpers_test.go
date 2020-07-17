// Copyright 2020 Google LLC
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

package server_test

import (
	"errors"
	"time"

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

type stubRawChainSource struct {
	data *certchain.RawChain
}

func (s *stubRawChainSource) Fetch(*certchain.RawChain, func() time.Time) (*certchain.RawChain, futureevent.Event, error) {
	return s.data, futureevent.NeverOccurs(), nil
}

type stubOCSPRespSource struct {
	data *certchain.OCSPResponse
}

func (s *stubOCSPRespSource) Fetch(*certchain.RawChain, func() time.Time) (*certchain.OCSPResponse, futureevent.Event, error) {
	return s.data, futureevent.NeverOccurs(), nil
}

//------------------------
//  waitResult

// waitResult represents the result of waitFor.
type waitResult int

const (
	waitSuccess waitResult = iota
	waitCanceled
	waitTimeout
)

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
	// TODO(banaag): unsupported now, will support with in-memory cache PR.
	return nil, certmanager.ErrNotFound
}

func (c *stubCache) Write(ac *certchain.AugmentedChain) error {
	if ac == nil {
		return errors.New("Write: nil augmented chain")
	}
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
