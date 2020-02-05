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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/processor/htmlproc/htmltask"
	"github.com/google/webpackager/resource/preload"
	"github.com/google/webpackager/resource/preload/preloadtest"
)

func TestPreloadStylesheets(t *testing.T) {
	pl := preloadtest.NewPreloadForRawLink

	tests := []struct {
		name string
		url  string
		html string
		want []*preload.Preload
	}{
		{
			name: "Simple",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <link rel="stylesheet" href="https://example.com/hello/foo.css">
			         <link rel="stylesheet" href="bar.css">
			       </head>`,
			want: []*preload.Preload{
				pl(`<https://example.com/hello/foo.css>;rel="preload";as="style"`),
				pl(`<https://example.com/hello/bar.css>;rel="preload";as="style"`),
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
			want: []*preload.Preload{
				pl(`<https://example.com/hello/foo.css>;rel="preload";as="style"`),
				pl(`<https://example.com/world/bar.css>;rel="preload";as="style"`),
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
			want: []*preload.Preload{
				pl(`<https://example.com/hello/baz.css>;rel="preload";as="style"`),
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
			want: nil,
		},
		{
			name: "CrossOrigin",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <link rel="stylesheet" href="https://example.org/hello/foo.css">
			       </head>`,
			want: []*preload.Preload{
				pl(`<https://example.org/hello/foo.css>;rel="preload";as="style"`),
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
			want: []*preload.Preload{
				pl(`<https://example.com/hello/foo.css>;rel="preload";as="style"`),
				pl(`<https://example.com/hello/bar.css>;rel="preload";as="style"`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := makeHTMLResponse(test.url, test.html)
			if err := htmltask.PreloadStylesheets().Run(resp); err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if diff := cmp.Diff(test.want, resp.Preloads); diff != "" {
				t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
