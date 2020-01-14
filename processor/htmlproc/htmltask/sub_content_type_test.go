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

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor/htmlproc/htmltask"
)

func TestExtractSubContentTypes(t *testing.T) {
	tests := []struct {
		name string
		url  string
		html string
		want []string
	}{
		{
			name: "CSS_Internal",
			url:  "https://example.com/hello/",
			html: `<!doctype html><style>body { font-family: serif; }</style>`,
			want: []string{"text/style"},
		},
		{
			name: "CSS_External",
			url:  "https://example.com/hello/",
			html: `<!doctype html><link rel="stylesheet" href="style.css">`,
			want: nil,
		},
		{
			name: "CSS_OtherMIMEType",
			url:  "https://example.com/hello/",
			html: `<!doctype html><style type="text/x-sass">/* ... */</style>`,
			want: []string{"text/x-sass"},
		},
		{
			name: "JS_Internal",
			url:  "https://example.com/hello/",
			html: `<!doctype html><script>(function () {/*...*/})();</script>`,
			want: []string{"application/javascript"},
		},
		{
			name: "JS_External",
			url:  "https://example.com/hello/",
			html: `<!doctype html><script src="script.js"></script>`,
			want: nil,
		},
		{
			name: "JS_OtherMIMEType",
			url:  "https://example.com/hello/",
			html: `<!doctype html><script type="text/plain">template</script>`,
			want: []string{"text/plain"},
		},
		{
			name: "SVG_Internal",
			url:  "https://example.com/hello/",
			html: `<!doctype html><svg><circle cx="10" cy="10" r="5" /></svg>`,
			want: []string{"image/svg+xml"},
		},
		{
			name: "SVG_External",
			url:  "https://example.com/hello/",
			html: `<!doctype html><img src="circle.svg" alt="circle">`,
			want: nil,
		},
		{
			name: "MathML",
			url:  "https://example.com/hello/",
			html: `<!doctype html><math><mi>x</mi><mo>+</mo><mn>1</mn></math>`,
			want: []string{"application/mathml+xml"},
		},
		{
			name: "Everything",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <style>body { font-family: serif; }</style>
			       <script>(function () {/*...*/})();</script>
			       <svg><circle cx="10" cy="10" r="5" /></svg>
			       <math><mi>x</mi><mo>+</mo><mn>1</mn></math>`,
			want: []string{
				"text/style",
				"application/javascript",
				"image/svg+xml",
				"application/mathml+xml",
			},
		},
	}

	task := htmltask.ExtractSubContentTypes()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := makeHTMLResponse(test.url, test.html)
			if err := task.Run(resp); err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if got := resp.ExtraData[exchange.SubContentType]; !reflect.DeepEqual(got, test.want) {
				t.Errorf("got = %#q, want %#q", got, test.want)
			}
		})
	}
}
