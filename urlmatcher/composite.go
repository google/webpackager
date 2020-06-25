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
)

// NotMatcher negates the Matcher M.
type NotMatcher struct{ M Matcher }

// Not returns a Matcher that negates m.
func Not(m Matcher) NotMatcher {
	return NotMatcher{m}
}

// Match returns true if neg.M.Match(u) returns false.
func (neg NotMatcher) Match(u *url.URL) bool {
	return !neg.M.Match(u)
}

// AndMatcher matches the URL u if u matches all of the given Matchers.
//
// An empty AndMatcher matches all URLs.
type AndMatcher []Matcher

// AllOf creates a new AndMatcher with m.
func AllOf(m ...Matcher) AndMatcher {
	return AndMatcher(m)
}

// Match calls Matchers in order and returns true if all of them return true.
// It stops the evaluation and returns false immediately when some Matcher
// returns false; the subsequent Matchers will not be called (short-circuit).
//
// Match returns true if am is empty.
func (am AndMatcher) Match(u *url.URL) bool {
	for _, m := range am {
		if !m.Match(u) {
			return false
		}
	}
	return true
}

// OrMatcher matches the URL u if u matches at least one of the given
// Matchers.
//
// An empty OrMatcher matches no URLs.
type OrMatcher []Matcher

// AnyOf creates a new OrMatcher with m.
func AnyOf(m ...Matcher) OrMatcher {
	return OrMatcher(m)
}

// Match calls Matchers in order and returns true if some Matcher returns true.
// It stops the evaluation immediately once it finds a matching; the subsequent
// Matchers will not be called (short-circuit).
//
// Match returns false if om is empty.
func (om OrMatcher) Match(u *url.URL) bool {
	for _, m := range om {
		if m.Match(u) {
			return true
		}
	}
	return false
}
