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

package urlmatcher

import (
	"net/url"
	"regexp"
	"strings"
)

// asciiEqualFold is like strings.EqualFold, but only folds ASCII letters
// to match the behavior with Section 3, RFC 4343.
func asciiEqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	// Use a C-style loop rather than a range loop to compare the strings
	// byte-by-byte. Note non-ASCII characters use only bytes of 0x80-0xFF
	// since Go uses UTF-8 encoding.
	for i := 0; i < len(s); i++ {
		if asciiToLower(s[i]) != asciiToLower(t[i]) {
			return false
		}
	}
	return true
}

func asciiToLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// HasScheme matches the URL u if u.Scheme is equal to scheme.
func HasScheme(scheme string) Matcher {
	return MatcherFunc(func(u *url.URL) bool {
		return u.Scheme == scheme
	})
}

// HasHost matches the URL u if u.Host is equal to host. ASCII letters are
// compared in a case-insensitive manner as specified in RFC 4343.
//
// Note u.Host may include a port number. To exclude it from matching, use
// HasHostname instead.
func HasHost(host string) Matcher {
	return MatcherFunc(func(u *url.URL) bool {
		return asciiEqualFold(u.Host, host)
	})
}

// HasHostname matches the URL u if u.Hostname() is equal to hostname.
// It does not inspect the port number, unlike HasHost. ASCII letters are
// compared in a case-insensitive manner as specified in RFC 4343.
func HasHostname(hostname string) Matcher {
	return MatcherFunc(func(u *url.URL) bool {
		return asciiEqualFold(u.Hostname(), hostname)
	})
}

// HasHostnameSuffix matches the URL u if u.Hostname() ends with suffix.
// If suffix has a leading dot ("."), the Matcher matches strictly the
// subdomains. Otherwise the Matcher matches the exact domain as well.
// In either case, the Matcher does not match any domains which are not
// a subdomain of suffix. For example:
//
//   - ".example.com" matches "www.example.com" but not "example.com".
//   - "example.com" matches both "www.example.com" and "example.com".
//   - Neither ".example.com" nor "exmple.com" matches "counterexample.com".
//
// ASCII letters are compared in a case-insensitive manner as specified in
// RFC 4343.
func HasHostnameSuffix(suffix string) Matcher {
	trimmed := strings.TrimPrefix(suffix, ".")

	return MatcherFunc(func(u *url.URL) bool {
		if len(u.Hostname()) < len(suffix) {
			return false
		}
		i := len(u.Hostname()) - len(trimmed)
		if i > 0 && u.Hostname()[i-1] != '.' {
			return false
		}
		return asciiEqualFold(u.Hostname()[i:], trimmed)
	})
}

// HasEscapedPathPrefix matches the URL u if u.EscapedPath() starts with
// prefix.
func HasEscapedPathPrefix(prefix string) Matcher {
	return MatcherFunc(func(u *url.URL) bool {
		return strings.HasPrefix(u.EscapedPath(), prefix)
	})
}

// HasEscapedPathRegexp matches the URL u if u.EscapedPath() matches the
// provided Regexp re. Note it just uses re.MatchString to determine the
// matching, so partial matches are considered as a match. If you want to
// match your pattern against the entire u.EscapedPath(), enclose it with
// `\A(?:...)\z`.
func HasEscapedPathRegexp(re *regexp.Regexp) Matcher {
	return MatcherFunc(func(u *url.URL) bool {
		return re.MatchString(u.EscapedPath())
	})
}

// HasRawQueryRegexp matches the URL u if u.RawQuery matches the provided
// Regexp re. Note it just uses re.MatchString to determine the matching,
// so partial matches are considered as a match. If you want to match your
// pattern against the entire u.RawQuery, enclose it with `\A(?:...)\z`.
func HasRawQueryRegexp(re *regexp.Regexp) Matcher {
	return MatcherFunc(func(u *url.URL) bool {
		return re.MatchString(u.RawQuery)
	})
}
