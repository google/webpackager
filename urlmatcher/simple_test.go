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
	"regexp"
	"testing"

	"github.com/google/webpackager/urlmatcher"
)

func TestHasScheme(t *testing.T) {
	matcher := urlmatcher.HasScheme("https")

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Match",
			url:  "https://example.com/hello/",
			want: true,
		},
		{
			name: "Mismatch",
			url:  "http://example.com/hello/",
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

func TestHasHost(t *testing.T) {
	matcher := urlmatcher.HasHost("example.com")

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "LowerCase",
			url:  "https://example.com/",
			want: true,
		},
		{
			name: "MixedCase",
			url:  "https://ExAmpLe.CoM/",
			want: true,
		},
		{
			name: "Subdomain",
			url:  "https://www.example.com/",
			want: false,
		},
		{
			name: "OtherDomain",
			url:  "https://example.org/",
			want: false,
		},
		{
			name: "ConfusingDomain",
			url:  "https://wwwexample.com/",
			want: false,
		},
		{
			name: "WithNonemptyPath",
			url:  "https://example.com/hello/",
			want: true,
		},
		{
			name: "WithPortNumber",
			url:  "https://example.com:10443/",
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

func TestHasHostname(t *testing.T) {
	matcher := urlmatcher.HasHostname("example.com")

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "LowerCase",
			url:  "https://example.com/",
			want: true,
		},
		{
			name: "MixedCase",
			url:  "https://ExAmpLe.Com/",
			want: true,
		},
		{
			name: "Subdomain",
			url:  "https://www.example.com/",
			want: false,
		},
		{
			name: "OtherDomain",
			url:  "https://example.org/",
			want: false,
		},
		{
			name: "ConfusingDomain",
			url:  "https://wwwexample.com/",
			want: false,
		},
		{
			name: "WithNonemptyPath",
			url:  "https://example.com/hello/",
			want: true,
		},
		{
			name: "WithPortNumber",
			url:  "https://example.com:10443/",
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

func TestHasHostnameSuffixWithoutDot(t *testing.T) {
	matcher := urlmatcher.HasHostnameSuffix("example.com")

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "LowerCase",
			url:  "https://example.com/",
			want: true,
		},
		{
			name: "MixedCase",
			url:  "https://ExAmpLe.Com/",
			want: true,
		},
		{
			name: "Subdomain",
			url:  "https://www.example.com/",
			want: true,
		},
		{
			name: "Subdomain_MixedCase",
			url:  "https://WWW.eXaMPLe.CoM/",
			want: true,
		},
		{
			name: "OtherDomain",
			url:  "https://example.org/",
			want: false,
		},
		{
			name: "ConfusingDomain",
			url:  "https://wwwexample.com/",
			want: false,
		},
		{
			name: "WithNonemptyPath",
			url:  "https://www.example.com/hello/",
			want: true,
		},
		{
			name: "WithPortNumber",
			url:  "https://www.example.com:10443/",
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

func TestHasHostnameSuffixWithDot(t *testing.T) {
	matcher := urlmatcher.HasHostnameSuffix(".example.com")

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "LowerCase",
			url:  "https://example.com/",
			want: false,
		},
		{
			name: "MixedCase",
			url:  "https://ExAmpLe.Com/",
			want: false,
		},
		{
			name: "Subdomain",
			url:  "https://www.example.com/",
			want: true,
		},
		{
			name: "Subdomain_MixedCase",
			url:  "https://WWW.eXaMPLe.CoM/",
			want: true,
		},
		{
			name: "OtherDomain",
			url:  "https://example.org/",
			want: false,
		},
		{
			name: "ConfusingDomain",
			url:  "https://wwwexample.com/",
			want: false,
		},
		{
			name: "WithNonemptyPath",
			url:  "https://www.example.com/hello/",
			want: true,
		},
		{
			name: "WithPortNumber",
			url:  "https://www.example.com:10443/",
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

func TestHasEscapedPathPrefix(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		url  string
		want bool
	}{
		{
			name: "ExactMatch",
			arg:  "/hello/",
			url:  "https://example.com/hello/",
			want: true,
		},
		{
			name: "PrefixMatch",
			arg:  "/hello/",
			url:  "https://example.com/hello/world/",
			want: true,
		},
		{
			name: "ObviousMismatch",
			arg:  "/foo/bar/baz/",
			url:  "https://example.com/hello/world/",
			want: false,
		},
		{
			name: "ArgNotURLEscaped",
			arg:  "/#/q=foo/",
			url:  "https://example.com/%23/q%3Dfoo/",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matcher := urlmatcher.HasEscapedPathPrefix(test.arg)
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := matcher.Match(u); got != test.want {
				t.Errorf("HasEscapedPathPrefix(%q).Match(%q) = %v, want %v", test.arg, u, got, test.want)
			}
		})
	}
}

func TestHasEscapedPathRegexp(t *testing.T) {
	tests := []struct {
		name string
		re   *regexp.Regexp
		url  string
		want bool
	}{
		{
			name: "Match",
			re:   regexp.MustCompile("^/hello/"),
			url:  "https://example.com/hello/world/",
			want: true,
		},
		{
			name: "NotMatch",
			re:   regexp.MustCompile("^/hello/"),
			url:  "https://example.com/foo/bar/baz/",
			want: false,
		},
		{
			name: "RegexpNotURLEscaped",
			re:   regexp.MustCompile("^/#/q=.*/$"),
			url:  "https://example.com/%23/q%3Dfoo/",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matcher := urlmatcher.HasEscapedPathRegexp(test.re)
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := matcher.Match(u); got != test.want {
				t.Errorf("HasEscapedPathRegexp(%q).Match(%q) = %v, want %v", test.re, u, got, test.want)
			}
		})
	}
}

func TestHasRawQueryRegexp(t *testing.T) {
	tests := []struct {
		name string
		re   *regexp.Regexp
		url  string
		want bool
	}{
		{
			name: "Match",
			re:   regexp.MustCompile("^q=[0-9]+$"),
			url:  "https://example.com/?q=42",
			want: true,
		},
		{
			name: "NotMatch",
			re:   regexp.MustCompile("^q=[0-9]+$"),
			url:  "https://example.com/?q=foo",
			want: false,
		},
		{
			name: "RegexpNotURLEscaped",
			re:   regexp.MustCompile("^q=#[^&]+$"),
			url:  "https://example.com/?q=%23foo",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matcher := urlmatcher.HasRawQueryRegexp(test.re)
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := matcher.Match(u); got != test.want {
				t.Errorf("HasRawQueryRegexp(%q).Match(%q) = %v, want %v", test.re, u, got, test.want)
			}
		})
	}
}
