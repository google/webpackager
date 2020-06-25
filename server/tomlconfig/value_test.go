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
	"time"
)

func TestParseJSExpiry(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  time.Duration
	}{
		{
			name:  "Simple_Middle",
			value: "12h",
			want:  12 * time.Hour,
		},
		{
			name:  "Simple_Maximum",
			value: "24h",
			want:  24 * time.Hour,
		},
		{
			name:  "Unsafe_Middle",
			value: "unsafe:100h",
			want:  100 * time.Hour,
		},
		{
			name:  "Unsafe_Maximum",
			value: "unsafe:168h",
			want:  168 * time.Hour,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseJSExpiry(test.value)
			if err != nil {
				t.Fatalf("parseJSExpiry(%q) = error(%q), want success", test.value, err)
			}
			if got != test.want {
				t.Fatalf("parseJSExpiry(%q) = %v, want %v", test.value, got, test.want)
			}
		})
	}
}

func TestParseJSExpiry_Error(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "Negative",
			value: "-1h",
		},
		{
			name:  "Zero",
			value: "0h",
		},
		{
			name:  "RequiresUnsafe",
			value: "100h",
		},
		{
			name:  "Exceeds7Days",
			value: "unsafe:169h",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseJSExpiry(test.value)
			if err == nil {
				t.Fatalf("parseJSExpiry(%q) = %v, want error", test.value, got)
			}
		})
	}
}

func TestMustCompileFullMatch(t *testing.T) {
	tests := []struct {
		name  string
		re    string
		value string
		want  bool
	}{
		{
			name:  "ExactMatch",
			re:    "/foo/bar/baz/",
			value: "/foo/bar/baz/",
			want:  true,
		},
		{
			name:  "RegexpMatch",
			re:    "/hello/.*",
			value: "/hello/world/",
			want:  true,
		},
		{
			name:  "ObviousMismatch",
			re:    "/foo/bar/baz/",
			value: "/hello/world/",
			want:  false,
		},
		{
			name:  "EmptyMatch",
			re:    "",
			value: "",
			want:  true,
		},
		{
			name:  "EmptyMismatch",
			re:    "",
			value: "/foo/bar/baz/",
			want:  false,
		},
		{
			name:  "MustMatchEntirely_NotPrefix",
			re:    "/foo/",
			value: "/foo/bar/baz/",
			want:  false,
		},
		{
			name:  "MustMatchEntirely_NotMiddle",
			re:    "/bar/",
			value: "/foo/bar/baz/",
			want:  false,
		},
		{
			name:  "MustMatchEntirely_NotSuffix",
			re:    "/baz/",
			value: "/foo/bar/baz/",
			want:  false,
		},
		{
			name:  "MustMatchEntirely_OrOperator",
			re:    "/foo|baz/",
			value: "/foo/bar/baz/",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := mustCompileFullMatch(test.re)
			if got := m.MatchString(test.value); got != test.want {
				t.Errorf("mustCompileFullMatch(%q).Match(%q) = %v, want %v", test.re, test.value, got, test.want)
			}
		})
	}
}
