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

func TestBoundedCache_LRU(t *testing.T) {
	foo := makeResource("https://example.com/foo.html")
	bar := makeResource("https://example.com/bar.html")
	bc := cache.NewBoundedInMemoryCache(1)

	reqFoo := makeRequest("https://example.com/foo.html")
	reqBar := makeRequest("https://example.com/bar.html")

	// foo is not present in the initial state.
	{
		got, err := bc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("bc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("bc.Lookup(reqFoo) = %v, want %v", got, nil)
		}
	}

	if err := bc.Store(foo); err != nil {
		t.Errorf("bc.Store(foo) = error(%q), want success", err)
	}

	// foo is now present.
	{
		got, err := bc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("bc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != foo {
			t.Errorf("bc.Lookup(reqFoo) = %v, want %v", got, foo)
		}
	}
	// bar is still not present.
	{
		got, err := bc.Lookup(reqBar)
		if err != nil {
			t.Errorf("bc.Lookup(reqBar) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("bc.Lookup(reqBar) = %v, want %v", got, nil)
		}
	}

	if err := bc.Store(bar); err != nil {
		t.Errorf("bc.Store(bar) = error(%q), want success", err)
	}

	// bar is now present.
	{
		got, err := bc.Lookup(reqBar)
		if err != nil {
			t.Errorf("bc.Lookup(reqBar) = error(%q), want success", err)
		}
		if got != bar {
			t.Errorf("bc.Lookup(reqBar) = %v, want %v", got, bar)
		}
	}
	// foo is no longer present.
	{
		got, err := bc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("bc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("bc.Lookup(reqFoo) = %v, want %v", got, nil)
		}
	}
}

func TestBoundedCache_2Q(t *testing.T) {
	foo := makeResource("https://example.com/foo.html")
	bar := makeResource("https://example.com/bar.html")
	baz := makeResource("https://example.com/baz.html")
	bc := cache.NewBoundedInMemoryCache(2)

	reqFoo := makeRequest("https://example.com/foo.html")
	reqBar := makeRequest("https://example.com/bar.html")
	reqBaz := makeRequest("https://example.com/baz.html")

	if err := bc.Store(foo); err != nil {
		t.Errorf("bc.Store(foo) = error(%q), want success", err)
	}
	if err := bc.Store(bar); err != nil {
		t.Errorf("bc.Store(bar) = error(%q), want success", err)
	}
	if err := bc.Store(baz); err != nil {
		t.Errorf("bc.Store(baz) = error(%q), want success", err)
	}

	// foo is not present.
	{
		got, err := bc.Lookup(reqFoo)
		if err != nil {
			t.Errorf("bc.Lookup(reqFoo) = error(%q), want success", err)
		}
		if got != nil {
			t.Errorf("bc.Lookup(reqFoo) = %v, want %v", got, nil)
		}
	}
	// bar is present.
	{
		got, err := bc.Lookup(reqBar)
		if err != nil {
			t.Errorf("bc.Lookup(reqBar) = error(%q), want success", err)
		}
		if got != bar {
			t.Errorf("bc.Lookup(reqBar) = %v, want %v", got, bar)
		}
	}
	// bar is present.
	{
		got, err := bc.Lookup(reqBaz)
		if err != nil {
			t.Errorf("bc.Lookup(reqBaz) = error(%q), want success", err)
		}
		if got != baz {
			t.Errorf("bc.Lookup(reqBaz) = %v, want %v", got, baz)
		}
	}
}
