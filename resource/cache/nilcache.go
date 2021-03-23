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
	"net/http"

	"github.com/google/webpackager/resource"
)

// NilCache returns a ResourceCache that stores and retrieves nothing.
func NilCache() ResourceCache {
	return globalNilCache
}

var globalNilCache = nilCache{}

type nilCache struct{}

func (_ nilCache) Lookup(req *http.Request) (*resource.Resource, error) {
	return nil, nil
}

func (_ nilCache) Store(r *resource.Resource) error {
	return nil
}
