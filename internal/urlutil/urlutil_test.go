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

package urlutil_test

import (
	"net/url"
	"testing"

	"github.com/google/webpackager/internal/urlutil"
)

func TestGetCleanPath(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "File",
			url:  "https://example.com/foo/index.html",
			want: "/foo/index.html",
		},
		{
			name: "File_NoExt",
			url:  "https://example.com/foo/bar",
			want: "/foo/bar",
		},
		{
			name: "Dir",
			url:  "https://example.com/foo/bar/",
			want: "/foo/bar/",
		},
		{
			name: "CurrentDir",
			url:  "https://example.com/foo/bar/.",
			want: "/foo/bar/",
		},
		{
			name: "ParentDir",
			url:  "https://example.com/foo/bar/..",
			want: "/foo/",
		},
		{
			name: "File_UncleahPath",
			url:  "https://example.com/foo/bar/../index.html",
			want: "/foo/index.html",
		},
		{
			name: "MultiSlashes",
			url:  "https://example.com/foo/bar//",
			want: "/foo/bar/",
		},
		{
			name: "Root",
			url:  "https://example.com/",
			want: "/",
		},
		{
			name: "Root_NoSlash",
			url:  "https://example.com",
			want: "/",
		},
		{
			name: "Root_Unclean",
			url:  "https://example.com/foo/bar/../..",
			want: "/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := urlutil.GetCleanPath(u); got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestHasSameOrigin(t *testing.T) {
	tests := []struct {
		name string
		u1   string
		u2   string
		want bool
	}{
		{
			name: "SameOrigin",
			u1:   "https://example.com/foo/",
			u2:   "https://example.com/bar/",
			want: true,
		},
		{
			name: "DifferentDomain",
			u1:   "https://example.org/foo/",
			u2:   "https://example.com/bar/",
			want: false,
		},
		{
			name: "DifferentScheme",
			u1:   "http://example.com/foo/",
			u2:   "https://example.com/bar/",
			want: false,
		},
		{
			name: "MissingScheme_u1",
			u1:   "//example.com/foo/",
			u2:   "https://example.com/bar/",
			want: false,
		},
		{
			name: "MissingScheme_u2",
			u1:   "https://example.com/foo/",
			u2:   "//example.com/bar/",
			want: false,
		},
		{
			name: "MissingScheme_both",
			u1:   "//example.com/foo/",
			u2:   "//example.com/bar/",
			want: false,
		},
		{
			name: "OpaqueOrigin",
			u1:   "data:image/png;base64,iVBORw0KGgo=",
			u2:   "data:image/png;base64,iVBORw0KGgo=",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u1, err := url.Parse(test.u1)
			if err != nil {
				t.Fatal(err)
			}
			u2, err := url.Parse(test.u2)
			if err != nil {
				t.Fatal(err)
			}
			if got := urlutil.HasSameOrigin(u1, u2); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "File",
			url:  "https://example.com/foo.html",
			want: false,
		},
		{
			name: "File_NoExt",
			url:  "https://example.com/foo",
			want: false,
		},
		{
			name: "Dir",
			url:  "https://example.com/foo/",
			want: true,
		},
		{
			name: "CurrentDir",
			url:  "https://example.com/foo/.",
			want: true,
		},
		{
			name: "ParentDir",
			url:  "https://example.com/foo/..",
			want: true,
		},
		{
			name: "File_UncleanPath",
			url:  "https://example.com/foo/../bar.html",
			want: false,
		},
		{
			name: "Root",
			url:  "https://example.com/",
			want: true,
		},
		{
			name: "Root_NoSlash",
			url:  "https://example.com",
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := urlutil.IsDir(u); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
