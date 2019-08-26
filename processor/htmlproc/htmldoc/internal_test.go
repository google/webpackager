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

package htmldoc

import (
	"net/url"
	"strings"
	"testing"

	"github.com/google/webpackager/internal/urlutil"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name string
		html string
		want *url.URL
	}{
		{
			name: "Absolute",
			html: `<!doctype html><base href="https://example.org/">`,
			want: urlutil.MustParse("https://example.org/"),
		},
		{
			name: "Relative",
			html: `<!doctype html><base href="foo/">`,
			want: urlutil.MustParse("foo/"),
		},
		{
			name: "ExtraAttribute",
			html: `<!doctype html><base target="_blank" href="foo/">`,
			want: urlutil.MustParse("foo/"),
		},
		{
			name: "Invalid",
			html: `<!doctype html><base href="INVALID :-)">`,
			want: &url.URL{},
		},
		{
			name: "Missing",
			html: `<!doctype html><base target="_blank">`,
			want: &url.URL{},
		},
		{
			name: "OtherElement",
			html: `<!doctype html><link rel="canonical" href="foo/">`,
			want: &url.URL{},
		},
		{
			name: "Misplaced",
			html: `<!doctype html><body><base href="foo/"></body>`,
			want: &url.URL{},
		},
		{
			name: "CompleteHTML",
			html: `
<!doctype html>
<html>
  <head>
    <meta charset="UTF-8">
    <base href="/pages/">
    <link href="style.css" rel="stylesheet">
  </head>
  <body>
    <div><a href="1.html">Sample 1</a></div>
    <div><a href="2.html">Sample 2</a></div>
  </body>
</html>`,
			want: urlutil.MustParse("/pages/"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(test.html))
			if err != nil {
				t.Fatal(err)
			}
			head := FindNode(doc, atom.Head)
			if head == nil {
				t.Fatal("head not found")
			}
			if got := getBaseURL(head); *got != *test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
