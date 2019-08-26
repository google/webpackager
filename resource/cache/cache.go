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

// Package cache defines the ResourceCache interface and provides the most
// basic implementation.
package cache

import (
	"net/http"

	"github.com/google/webpackager/resource"
)

// ResourceCache caches resource.Resources.
//
// ResourceCache implementations should match Resources in a way analogous
// to HTTP caches: the returned Resource should have a matching RequestURL
// and, if applicable, compatible HTTP Variants headers (but see Bugs with
// OnMemoryCache).
//
// For HTTP Variants, please see:
// https://httpwg.org/http-extensions/draft-ietf-httpbis-variants.html.
type ResourceCache interface {
	// Lookup returns a Resource matching req. It returns nil (with no error)
	// when ResourceCache contains no matching Resource.
	Lookup(req *http.Request) (*resource.Resource, error)

	// Store stores the provided Resource r into the cache.
	Store(r *resource.Resource) error
}
