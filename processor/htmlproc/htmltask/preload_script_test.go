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
)

func TestScriptTask(t *testing.T) {
	tests := []struct {
		name string
		url  string
		html string
		want []string
	}{
		{
			name: "Head",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <script src="https://example.com/hello/foo.js"></script>
			         <script src="bar.js"></script>
			         <title>Test Docuemnt</title>
			         <script src="baz.js"></script>
			       </head>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
				`<https://example.com/hello/bar.js>;rel="preload";as="script"`,
				`<https://example.com/hello/baz.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Body_PreceedsNone",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Body_Text",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         This is a test document.
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{},
		},
		{
			name: "Body_Image",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <img src="thumb.jpg">
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{},
		},
		{
			name: "Body_EmptyElement",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <div></div>
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Body_InnerText",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <div>This is a test document.</div>
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{},
		},
		{
			name: "Body_DirectScript",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
                     <script>(function(){})();</script>
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Body_Noscript",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
                     <noscript>Needs JavaScript.</noscript>
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Body_Comment",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <!-- Comment -->
			         <script src="foo.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Body_Complex",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <body>
			         <script src="foo.js"></script>
			         <div> <!-- An empty div. --> </div>
			         <script src="bar.js"></script>
			         <div>This is a test document.</div>
			         <script src="baz.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
				`<https://example.com/hello/bar.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "BaseURL",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <base href="/world/">
			         <script src="https://example.com/hello/foo.js"></script>
			         <script src="bar.js"></script>
			       </head>
			       <body>
			         <script src="baz.js"></script>
			       </body>`,
			want: []string{
				`<https://example.com/hello/foo.js>;rel="preload";as="script"`,
				`<https://example.com/world/bar.js>;rel="preload";as="script"`,
				`<https://example.com/world/baz.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "CrossDomain",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <script src="https://example.org/hello/foo.js"></script>
			       </head>`,
			want: []string{
				`<https://example.org/hello/foo.js>;rel="preload";as="script"`,
			},
		},
		{
			name: "Attribute",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <head>
			         <script src="async.js" async></script>
			         <script src="defer.js" defer></script>
			         <script src="type.js" type="text/javascript"></script>
			         <script>src="embed.js"</script>
			       </head>`,
			want: []string{
				`<https://example.com/hello/type.js>;rel="preload";as="script"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := makeHTMLResponse(test.url, test.html)
			if err := htmltask.InsecurePreloadScripts().Run(resp); err != nil {
				t.Errorf("got error(%q), want success", err)
			}
			if diff := cmp.Diff(test.want, preloadHeaders(resp)); diff != "" {
				t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
