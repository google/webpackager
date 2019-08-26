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

package cache_test

import (
	"net/http"
	"testing"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/cache"
)

func makeRequest(rawurl string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func makeResource(rawurl string) *resource.Resource {
	return resource.NewResource(urlutil.MustParse(rawurl))
}

func TestOnMemoryCache(t *testing.T) {
	foo := makeResource("https://example.com/foo.html")
	bar := makeResource("https://example.com/bar.html")
	mc := cache.NewOnMemoryCache()

	reqFoo := makeRequest("https://example.com/foo.html")
	reqBar := makeRequest("https://example.com/bar.html")

	// foo is not present in the initial state.
	{
		got, err := mc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("mc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("mc.Lookup(reqFoo) = %v, want %v", got, nil)
		}
	}

	if err := mc.Store(foo); err != nil {
		t.Errorf("mc.Store(foo) = error(%q), want success", err)
	}

	// foo is now present.
	{
		got, err := mc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("mc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != foo {
			t.Errorf("mc.Lookup(reqFoo) = %v, want %v", got, foo)
		}
	}
	// bar is still not present.
	{
		got, err := mc.Lookup(reqBar)
		if err != nil {
			t.Errorf("mc.Lookup(reqBar) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("mc.Lookup(reqBar) = %v, want %v", got, nil)
		}
	}

	if err := mc.Store(bar); err != nil {
		t.Errorf("mc.Store(bar) = error(%q), want success", err)
	}
}
