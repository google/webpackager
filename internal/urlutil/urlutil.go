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

// Package urlutil provides URL-related utility functions.
package urlutil

import (
	"net/url"
	"path"
)

// GetCleanPath returns cleaned u.Path. It is like path.Clean(u.Path) but does
// not remove the trailing slash.
func GetCleanPath(u *url.URL) string {
	p := path.Clean(u.Path)
	if p == "." || p == "/" {
		return "/"
	}
	if IsDir(u) {
		p += "/"
	}
	return p
}

// HasSameOrigin reports whether u1 and u2 have the same nonempty Scheme and
// Host. u1 and u2 are assumed to be absolute URLs with the http/https scheme.
// Nevertheless, HasSameOrigin returns false when either or both of them have
// empty Scheme or Host.
//
// HasSameOrigin does not implement all the rules defined for "same origin"
// in [HTML5] but yields the same result if u1 and u2 are both absolute URLs
// with the http/https scheme, barring one corner case: HasSameOrigin does not
// recognize the default port numbers, e.g. it reports https://example.com and
// https://example.com:443 have a "different origin."
//
// [HTML5]: https://www.w3.org/TR/html52/browser.html#section-origin
func HasSameOrigin(u1 *url.URL, u2 *url.URL) bool {
	if u1.Scheme == "" || u1.Host == "" {
		return false
	}
	return u1.Scheme == u2.Scheme && u1.Host == u2.Host
}

// IsDir reports whether u represents a directory. Specifically, IsDir returns
// true when u.Path ends with a slash ("/") or the last component of u.Path is
// either "." or "..".
func IsDir(u *url.URL) bool {
	_, last := path.Split(u.Path)
	return last == "" || last == "." || last == ".."
}

// MustParse is like url.Parse but panics on a parse error. It simplifies safe
// initialization of global variables holding parsed URLs.
func MustParse(rawurl string) *url.URL {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return url
}
