// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"fmt"
	"net/http"

	"github.com/google/webpackager/resource"

	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/xerrors"
)

// NewBoundedInMemoryCache returns a new ResourceCache that stores Resources in
// memory, with an eviction policy after `size` entries. size must be positive.
func NewBoundedInMemoryCache(size int) ResourceCache {
	if size > 1 {
		// The extra memory/CPU overhead of lru.TwoQueueCache over lru.Cache
		// seems worth it to increase cache hit rate, given the comparatively
		// high expense to regenerate entries.
		c, err := lru.New2Q(size)
		if err != nil {
			// Only occurs if size < 2.
			panic(xerrors.Errorf("constructing 2Q cache: %w", err))
		}
		return &boundedCache{c}
	} else {
		// lru.New2Q can't construct a TwoQueueCache of size 1, because
		// it rounds 1*ratio down to 0 when constructing its inner
		// recent and recentEvict queues. See
		// https://github.com/hashicorp/golang-lru/issues/51.
		c, err := lru.New(size)
		if err != nil {
			// Only occurs if size < 1.
			panic(xerrors.Errorf("constructing LRU cache: %w", err))
		}
		return &boundedCache{lruCache{c}}
	}
}

type boundedCache struct {
	cache cache
}

func (c *boundedCache) Lookup(req *http.Request) (*resource.Resource, error) {
	switch r, _ := c.cache.Get(req.URL.String()); t := r.(type) {
	case *resource.Resource:
		return t, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected cache entry type %T", t)
	}
}

func (c *boundedCache) Store(r *resource.Resource) error {
	c.cache.Add(r.RequestURL.String(), r)
	return nil
}

type cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

// Compiler check that both lruCache and lru.TwoQueueCache implement cache:
var (
	_ cache = lruCache{}
	_ cache = (*lru.TwoQueueCache)(nil)
)

// Wrapper for *lru.Cache that elides the return value from Add, to match
// lru.TwoQueueCache.
type lruCache struct {
	lru *lru.Cache
}

func (c lruCache) Add(key, value interface{}) {
	c.lru.Add(key, value)
}

func (c lruCache) Get(key interface{}) (value interface{}, ok bool) {
	return c.lru.Get(key)
}
