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

package urlrewrite_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/webpackager/urlrewrite"
)

func TestRuleSequence(t *testing.T) {
	rs := urlrewrite.RuleSequence{
		&appendToPath{"1/"},
		&appendToPath{"2/"},
		&appendToPath{"3/"},
	}

	u, err := url.Parse("https://example.com/")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://example.com/1/2/3/"
	rs.Rewrite(u, http.Header{})
	if u.String() != want {
		t.Errorf("got %q, want %q", u, want)
	}
}

type appendToPath struct {
	s string
}

func (rule *appendToPath) Rewrite(u *url.URL, respHeader http.Header) {
	u.Path += rule.s
}
