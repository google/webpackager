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

package fetch_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/webpackager/fetch"
)

func newGetRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

type failingTweaker struct {
	err error
}

func (t *failingTweaker) Tweak(req, parent *http.Request) error {
	return t.err
}

func TestRequestTweakerSequence(t *testing.T) {
	errDummy := errors.New("errDummy")

	parent := newGetRequest("https://example.com/index.html")

	tests := []struct {
		name    string
		tweaker fetch.RequestTweaker
		want    http.Header
		err     error
	}{
		{
			name: "Success",
			tweaker: fetch.RequestTweakerSequence{
				fetch.SetCustomHeaders(http.Header{"X-Test-1": []string{"foo"}}),
				fetch.SetCustomHeaders(http.Header{"X-Test-2": []string{"bar"}}),
				fetch.SetReferer(),
			},
			want: http.Header{
				"X-Test-1": []string{"foo"},
				"X-Test-2": []string{"bar"},
				"Referer":  []string{"https://example.com/index.html"},
			},
			err: nil,
		},
		{
			name: "Failure",
			tweaker: fetch.RequestTweakerSequence{
				fetch.SetCustomHeaders(http.Header{"X-Test-1": []string{"foo"}}),
				&failingTweaker{errDummy},
				fetch.SetCustomHeaders(http.Header{"X-Test-2": []string{"bar"}}),
				fetch.SetReferer(),
			},
			want: http.Header{
				"X-Test-1": []string{"foo"},
			},
			err: errDummy,
		},
	}

	for _, test := range tests {
		req := newGetRequest("https://example.com/style.css")

		if err := test.tweaker.Tweak(req, parent); err != test.err {
			t.Errorf("got %v, want %v", err, test.err)
		}
		if !reflect.DeepEqual(req.Header, test.want) {
			t.Errorf("req.Header = %q, want %q", req.Header, test.want)
		}
	}
}

func TestSetReferer(t *testing.T) {
	tests := []struct {
		name   string
		parent *http.Request
		want   string
	}{
		{
			name:   "WithParent",
			parent: newGetRequest("https://example.com/index.html"),
			want:   "https://example.com/index.html",
		},
		{
			name:   "WithoutParent",
			parent: nil,
			want:   "",
		},
	}

	setReferer := fetch.SetReferer()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := newGetRequest("https://example.com/style.css")
			if err := setReferer.Tweak(req, test.parent); err != nil {
				t.Fatal(err)
			}
			if got := req.Referer(); got != test.want {
				t.Errorf("req.Referer() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestCopyParentHeaders(t *testing.T) {
	tests := []struct {
		name    string
		parent  http.Header
		before  http.Header
		tweaker fetch.RequestTweaker
		after   http.Header
	}{
		{
			name: "BaseCase",
			parent: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
				"From":            []string{"admin@example.com"},
			},
			before: http.Header{},
			tweaker: fetch.CopyParentHeaders([]string{
				"Accept-Encoding",
				"Accept-Language",
			}),
			after: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
			},
		},
		{
			name: "NonCanonicalKeys",
			parent: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
				"From":            []string{"admin@example.com"},
			},
			before: http.Header{},
			tweaker: fetch.CopyParentHeaders([]string{
				"accept-encoding",
				"accept-language",
			}),
			after: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
			},
		},
		{
			name: "CopyImplyOverwrite",
			parent: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
				"From":            []string{"admin@example.com"},
			},
			before: http.Header{
				"Accept-Language": []string{"ja-JP", "ja;q=0.9"},
				"From":            []string{"user@example.org"},
			},
			tweaker: fetch.CopyParentHeaders([]string{
				"Accept-Encoding",
				"Accept-Language",
			}),
			after: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
				"From":            []string{"user@example.org"},
			},
		},
		{
			name: "FieldsMissing",
			parent: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"From":            []string{"admin@example.com"},
			},
			before: http.Header{
				"Accept-Language": []string{"ja-JP", "ja;q=0.9"},
				"From":            []string{"user@example.org"},
			},
			tweaker: fetch.CopyParentHeaders([]string{
				"Accept-Encoding",
				"Accept-Language",
			}),
			after: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"ja-JP", "ja;q=0.9"},
				"From":            []string{"user@example.org"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			par := newGetRequest("https://example.com/index.html")
			par.Header = test.parent

			req := newGetRequest("https://example.com/style.css")
			req.Header = test.before

			if err := test.tweaker.Tweak(req, par); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(req.Header, test.after) {
				t.Errorf("req.Header = %q, want %q", req.Header, test.after)
			}
		})
	}
}

func TestCopyParentHeaders_SliceCloned(t *testing.T) {
	par := newGetRequest("https://example.com/index.html")
	req := newGetRequest("https://example.com/style.css")

	userAgent := "request_test/0.1"
	par.Header.Add("User-Agent", userAgent)

	tweaker := fetch.CopyParentHeaders([]string{"User-Agent"})
	if err := tweaker.Tweak(req, par); err != nil {
		t.Fatal(err)
	}
	req.Header["User-Agent"][0] = "(overwritten)"
	if got := par.Header.Get("User-Agent"); got != userAgent {
		t.Errorf(`par.Header["User-Agent"] = %q, want %q`, got, userAgent)
	}
}

func TestSetCustomHeaders(t *testing.T) {
	tests := []struct {
		name    string
		before  http.Header
		tweaker fetch.RequestTweaker
		after   http.Header
	}{
		{
			name:   "BaseCase",
			before: http.Header{},
			tweaker: fetch.SetCustomHeaders(http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
			}),
			after: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
			},
		},
		{
			name: "CopyImplyOverwrite",
			before: http.Header{
				"Accept-Language": []string{"ja-JP", "ja;q=0.9"},
				"User-Agent":      []string{"request_test/0.1"},
			},
			tweaker: fetch.SetCustomHeaders(http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
			}),
			after: http.Header{
				"Accept-Encoding": []string{"gzip"},
				"Accept-Language": []string{"en-US", "en;q=0.9"},
				"User-Agent":      []string{"request_test/0.1"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := newGetRequest("https://example.com/style.css")
			req.Header = test.before

			if err := test.tweaker.Tweak(req, nil); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(req.Header, test.after) {
				t.Errorf("req.Header = %q, want %q", req.Header, test.after)
			}
		})
	}
}

func TestSetCustomHeaders_SliceCloned(t *testing.T) {
	req := newGetRequest("https://example.com/style.css")

	userAgent := "request_test/0.1"
	header := http.Header{"User-Agent": []string{userAgent}}

	tweaker := fetch.SetCustomHeaders(header)
	if err := tweaker.Tweak(req, nil); err != nil {
		t.Fatal(err)
	}
	req.Header["User-Agent"][0] = "(overwritten)"
	if got := header.Get("User-Agent"); got != userAgent {
		t.Errorf(`header["User-Agent"] = %q, want %q`, got, userAgent)
	}
}
