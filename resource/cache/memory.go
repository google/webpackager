// Copyright 2019 Google LLC
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
	"net/http"

	"github.com/google/webpackager/resource"
)

// NewOnMemoryCache creates and initializes a new ResourceCache storing
// Resources on memory.
func NewOnMemoryCache() ResourceCache {
	return make(onMemoryCache)
}

// BUG(yuizumi): OnMemoryCache uses only RequestURL for the cache key at this
// moment; it is not aware of Vary or Variants yet.
type onMemoryCache map[string]*resource.Resource

func (mc onMemoryCache) Lookup(req *http.Request) (*resource.Resource, error) {
	return mc[req.URL.String()], nil
}

func (mc onMemoryCache) Store(r *resource.Resource) error {
	mc[r.RequestURL.String()] = r
	return nil
}
