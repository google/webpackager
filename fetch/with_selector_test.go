// Copyright 2020 Google LLC
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

package fetch_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/urlmatcher"
)

func TestSelector(t *testing.T) {
	selector := &fetch.Selector{
		Allow: []urlmatcher.Matcher{
			urlmatcher.HasHostnameSuffix("example.com"),
			urlmatcher.HasHostnameSuffix("example.org"),
		},
		Deny: []urlmatcher.Matcher{
			urlmatcher.HasHostname("bad.example.com"),
		},
	}

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Allowed_1stRule",
			url:  "https://www.example.com/",
			want: true,
		},
		{
			name: "Allowed_2ndRule",
			url:  "https://www.example.org/",
			want: true,
		},
		{
			name: "DeniedExplicitly",
			url:  "https://bad.example.com/",
			want: false,
		},
		{
			name: "NotAllowed",
			url:  "https://www.example.net/",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := selector.Match(u); got != test.want {
				t.Errorf("Match(%q) = %v, want %v", u, got, test.want)
			}
		})
	}
}

func TestSelector_EmptySlices(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		selector *fetch.Selector
		want     bool
	}{
		{
			name: "EmptyAllow_MeansAllowAny",
			url:  "https://www.example.org/",
			selector: &fetch.Selector{
				Allow: nil,
				Deny: []urlmatcher.Matcher{
					urlmatcher.HasHostname("bad.example.com"),
					urlmatcher.HasHostname("bad.example.org"),
				},
			},
			want: true,
		},
		{
			name: "EmptyAllow_DenyStillApplies",
			url:  "https://bad.example.org/",
			selector: &fetch.Selector{
				Allow: nil,
				Deny: []urlmatcher.Matcher{
					urlmatcher.HasHostname("bad.example.com"),
					urlmatcher.HasHostname("bad.example.org"),
				},
			},
			want: false,
		},
		{
			name: "EmptyDeny_MeansDenyNone",
			url:  "https://www.example.com/",
			selector: &fetch.Selector{
				Allow: []urlmatcher.Matcher{
					urlmatcher.HasHostnameSuffix("example.com"),
				},
				Deny: nil,
			},
			want: true,
		},
		{
			name: "EmptyDeny_AllowStillApplies",
			url:  "https://www.example.org/",
			selector: &fetch.Selector{
				Allow: []urlmatcher.Matcher{
					urlmatcher.HasHostnameSuffix("example.com"),
				},
				Deny: nil,
			},
			want: false,
		},
		{
			name: "EmptyBoth_MeansAllowAnyDenyNone",
			url:  "https://www.example.com/",
			selector: &fetch.Selector{
				Allow: nil,
				Deny:  nil,
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := test.selector.Match(u); got != test.want {
				t.Errorf("Match(%q) = %v, want %v", u, got, test.want)
			}
		})
	}
}

func TestWithSelector_URLMatch(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	f := fetch.WithSelector(&stubFetcher{}, urlmatcher.HasHost("example.com"))
	resp, err := f.Do(req)
	if err != nil {
		t.Fatalf("Fetch() = error(%q), want success", err)
	}
	defer resp.Body.Close()
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]byte("<!doctype html><p>hello</p>"), got); diff != "" {
		t.Errorf("Fetch() mismatch (-want +got):\n%s", diff)
	}
}

func TestWithSelector_URLMismatch(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://example.org/", nil)
	if err != nil {
		t.Fatal(err)
	}
	f := fetch.WithSelector(&stubFetcher{}, urlmatcher.HasHost("example.com"))
	resp, err := f.Do(req)
	if err != fetch.ErrURLMismatch {
		if err != nil {
			t.Fatalf("Fetch() = error(%q), want ErrURLMismatch", err)
		} else {
			t.Errorf("Fetch() = %#v, want ErrURLMismatch", resp)
		}
	}
}
