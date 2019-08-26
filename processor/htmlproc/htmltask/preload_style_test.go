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

package htmltask_test

import (
	"reflect"
	"testing"

	"github.com/google/webpackager/processor/htmlproc/htmltask"
)

func TestPreloadStylesheets(t *testing.T) {
	tests := []struct {
		name string
		url  string
		html string
		want []string
	}{
		{
			name: "Simple",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <link rel="stylesheet" href="https://example.com/hello/foo.css">
			         <link rel="stylesheet" href="bar.css">
			       </head>`,
			want: []string{
				`<https://example.com/hello/foo.css>;rel="preload";as="style"`,
				`<https://example.com/hello/bar.css>;rel="preload";as="style"`,
			},
		},
		{
			name: "BaseURL",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <base href="/world/">
			         <link rel="stylesheet" href="https://example.com/hello/foo.css">
			         <link rel="stylesheet" href="bar.css">
			       </head>`,
			want: []string{
				`<https://example.com/hello/foo.css>;rel="preload";as="style"`,
				`<https://example.com/world/bar.css>;rel="preload";as="style"`,
			},
		},
		{
			name: "Alternate",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <link href="foo.css" title="foo" rel="alternate stylesheet">
			         <link href="bar.css" title="bar" rel="stylesheet alternate">
			         <link href="baz.css" title="baz" rel="stylesheet">
			       </head>`,
			want: []string{
				`<https://example.com/hello/baz.css>;rel="preload";as="style"`,
			},
		},
		{
			name: "OutsideHead",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <link rel="stylesheet" href="https://example.com/hello/foo.css">
			         <link rel="stylesheet" href="bar.css">
			       </body>`,
			want: []string{},
		},
		{
			name: "CrossOrigin",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <link rel="stylesheet" href="https://example.org/hello/foo.css">
			       </head>`,
			want: []string{
				`<https://example.org/hello/foo.css>;rel="preload";as="style"`,
			},
		},
		{
			name: "Tricky",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <link rel="test stylesheet" href="foo.css">
			       <!-- &#10; == LF ("\n") -->
			       <link rel="stylesheet&#10;" href="bar.css">
			       <link rel="stylesheet">`,
			want: []string{
				`<https://example.com/hello/foo.css>;rel="preload";as="style"`,
				`<https://example.com/hello/bar.css>;rel="preload";as="style"`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := makeHTMLResponse(test.url, test.html)
			if err := htmltask.PreloadStylesheets().Run(resp); err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if got := preloadHeaders(resp); !reflect.DeepEqual(got, test.want) {
				t.Errorf("resp.Preloads = %#q, want %#q", got, test.want)
			}
		})
	}
}
