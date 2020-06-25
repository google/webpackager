// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tomlconfig

import (
	"testing"
)

func TestVerifyServePath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "SingleSlash",
			path: "/",
		},
		{
			name: "WithoutTrailingSlash",
			path: "/foo/bar",
		},
		{
			name: "WithTrailingSlash",
			path: "/foo/bar/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := verifyServePath(test.path); err != nil {
				t.Errorf("verifyServePath(%q) = error(%q), want success", test.path, err)
			}
		})
	}
}

func TestVerifyServePath_Error(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "SingleDot",
			path: "/foo/./bar",
		},
		{
			name: "DoubleDot",
			path: "/foo/../bar",
		},
		{
			name: "NotAbsolute",
			path: "foo/bar",
		},
		{
			name: "DoubleSlash_Prefix",
			path: "//foo/bar",
		},
		{
			name: "DoubleSlash_Middle",
			path: "/foo//bar",
		},
		{
			name: "DoubleSlash_Suffix",
			path: "/foo/bar//",
		},
		{
			name: "DoubleSlash_Only",
			path: "//",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := verifyServePath(test.path); err == nil {
				t.Errorf("verifyServePath(%q) = success, want error", test.path)
			}
		})
	}
}

func TestVerifyURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "WithTrailingSlash_Path",
			url:  "/priv/doc/",
		},
		{
			name: "WithTrailingSlash_Full",
			url:  "https://example.com/priv/doc/",
		},
		{
			name: "WithoutTrailingSlash_Path",
			url:  "/priv/doc",
		},
		{
			name: "WithoutTrailingSlash_Full",
			url:  "https://example.com/priv/doc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := verifyURL(test.url); err != nil {
				t.Errorf("verifyURL(%q) = error(%q), want success", test.url, err)
			}
		})
	}
}

func TestVerifyURL_Error(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "SingleDot_Path",
			url:  "/priv/./doc",
		},
		{
			name: "SingleDot_Full",
			url:  "https://example.com/priv/./doc",
		},
		{
			name: "DoubleDot_Path",
			url:  "/priv/../doc",
		},
		{
			name: "DoubleDot_Full",
			url:  "https://example.com/priv/../doc",
		},
		{
			name: "DoubleSlash_Path",
			url:  "/priv//doc",
		},
		{
			name: "DoubleSlah_Full",
			url:  "https://example.com/priv//doc",
		},
		{
			name: "NotAbsolute",
			url:  "priv/doc",
		},
		{
			name: "InvalidScheme",
			url:  "http://example.com/priv/doc",
		},
		{
			name: "MissingScheme",
			url:  "//example.com/priv/doc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := verifyURL(test.url); err == nil {
				t.Errorf("verifyURL(%q) = success, want error", test.url)
			}
		})
	}
}
