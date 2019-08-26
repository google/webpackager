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

package urlrewrite

import (
	"net/http"
	"net/url"
)

// Rule implements a URL rewrite rule.
type Rule interface {
	// Rewrite mutates the provided URL u, applying the rewrite rule.
	// respHeader is the response header and used, for example, to replace
	// the file extension based on Content-Type.
	Rewrite(u *url.URL, respHeader http.Header)
}

// RuleSequence represents a series of rewrite rules.
type RuleSequence []Rule

// Rewrite applies all rewrite rules in order.
func (rs RuleSequence) Rewrite(u *url.URL, respHeader http.Header) {
	for _, r := range rs {
		r.Rewrite(u, respHeader)
	}
}
