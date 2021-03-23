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

package cache_test

import (
	"testing"

	"github.com/google/webpackager/resource/cache"
)

func TestNilCache(t *testing.T) {
	foo := makeResource("https://example.com/foo.html")
	nc := cache.NilCache()

	reqFoo := makeRequest("https://example.com/foo.html")

	// foo is not present in the initial state.
	{
		got, err := nc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("nc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("nc.Lookup(reqFoo) = %v, want %v", got, nil)
		}
	}

	if err := nc.Store(foo); err != nil {
		t.Errorf("nc.Store(foo) = error(%q), want success", err)
	}

	// foo is still not present.
	{
		got, err := nc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("nc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("nc.Lookup(reqFoo) = %v, want %v", got, nil)
		}
	}
}
