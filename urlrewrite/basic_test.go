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

func TestBasic(t *testing.T) {
	tests := []struct {
		name string
		url  string
		rule urlrewrite.Rule
		want string
	}{
		{
			name: "CleanPath",
			url:  "https://example.com/hello/../world/.",
			rule: urlrewrite.CleanPath(),
			want: "https://example.com/world/",
		},
		{
			name: "IndexRule_NotApplicable",
			url:  "https://example.com/foo/bar/hello.html",
			rule: urlrewrite.IndexRule("index.html"),
			want: "https://example.com/foo/bar/hello.html",
		},
		{
			name: "IndexRule_Applicable",
			url:  "https://example.com/foo/bar/",
			rule: urlrewrite.IndexRule("index.html"),
			want: "https://example.com/foo/bar/index.html",
		},
		{
			name: "IndexRule_WindowsServer",
			url:  "https://example.com/foo/bar/",
			rule: urlrewrite.IndexRule("default.asp"),
			want: "https://example.com/foo/bar/default.asp",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}
			test.rule.Rewrite(u, http.Header{})
			if u.String() != test.want {
				t.Errorf("got %q, want %q", u, test.want)
			}
		})
	}
}
