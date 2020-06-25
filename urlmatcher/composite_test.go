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

package urlmatcher_test

import (
	"net/url"
	"testing"

	"github.com/google/webpackager/urlmatcher"
)

func TestNot(t *testing.T) {
	matcher := urlmatcher.Not(urlmatcher.HasHost("example.com"))

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "ExampleCom",
			url:  "https://example.com/",
			want: false,
		},
		{
			name: "NotExampleCom",
			url:  "https://example.org/",
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := matcher.Match(u); got != test.want {
				t.Errorf("Match(%q) = %v, want %v", u, got, test.want)
			}
		})
	}
}

func TestAllOf(t *testing.T) {
	matcher := urlmatcher.AllOf(
		urlmatcher.HasHost("example.com"),
		urlmatcher.HasScheme("https"),
	)

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Match",
			url:  "https://example.com/",
			want: true,
		},
		{
			name: "Mismatch_Scheme",
			url:  "http://example.com/",
			want: false,
		},
		{
			name: "Mismatch_Host",
			url:  "https://example.org/",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := matcher.Match(u); got != test.want {
				t.Errorf("Match(%q) = %v, want %v", u, got, test.want)
			}
		})
	}
}

func TestAnyOf(t *testing.T) {
	matcher := urlmatcher.AnyOf(
		urlmatcher.HasHost("example.com"),
		urlmatcher.HasHost("example.org"),
	)

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Match_1st",
			url:  "https://example.com/",
			want: true,
		},
		{
			name: "Match_2nd",
			url:  "https://example.org/",
			want: true,
		},
		{
			name: "NoMatch",
			url:  "https://example.net/",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := matcher.Match(u); got != test.want {
				t.Errorf("Match(%q) = %v, want %v", u, got, test.want)
			}
		})
	}
}
