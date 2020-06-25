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

package fetch

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/google/webpackager/urlmatcher"
)

// ErrURLMismatch is returned by WithSelector clients when the request URL
// does not match their selector.
var ErrURLMismatch = errors.New("fetch: URL doesn't match the fetch targets")

// WithSelector wraps client to issue a fetch only if the request URL matches
// selector.
func WithSelector(client FetchClient, selector urlmatcher.Matcher) FetchClient {
	return &withSelector{client, selector}
}

type withSelector struct {
	client   FetchClient
	selector urlmatcher.Matcher
}

func (w *withSelector) Do(req *http.Request) (*http.Response, error) {
	if !w.selector.Match(req.URL) {
		return nil, ErrURLMismatch
	}
	return w.client.Do(req)
}

// Selector is a urlmatcher.Matcher designed for use with WithSelector. It is
// almost equivalent to
//
//     AllOf(AnyOf(Allow...), Not(AnyOf(Deny...)))
//
// but empty Allow is interpreted as "allow any." Empty Deny is interpreted
// as "deny none," like AnyOf(empty). Using Selector with WithSelector makes
// the selector easier to understand, especially for those familiar with the
// Allow/Deny style. Using Selector is not mandatory though: WithSelector can
// use any urlmatcher.Matcher as the selector.
type Selector struct {
	Allow []urlmatcher.Matcher
	Deny  []urlmatcher.Matcher
}

// Match implements the urlmatcher.Matcher interface.
func (s *Selector) Match(u *url.URL) bool {
	if len(s.Allow) > 0 && !urlmatcher.OrMatcher(s.Allow).Match(u) {
		return false
	}
	return !urlmatcher.OrMatcher(s.Deny).Match(u)
}
