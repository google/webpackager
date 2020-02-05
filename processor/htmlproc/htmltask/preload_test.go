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

func TestExtractPreloadTags(t *testing.T) {
	pl := preloadtest.NewPreloadForRawLink

	tests := []struct {
		name string
		url  string
		html string
		want []*preload.Preload
	}{
		{
			name: "Minimal",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <link href="https://example.com/hello/foo.jpg"
			           rel="preload" as="image">
			       <link href="bar.jpg" rel="preload" as="image">`,
			want: []*preload.Preload{
				pl(`<https://example.com/hello/foo.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/hello/bar.jpg>;rel="preload";as="image"`),
			},
		},
		{
			name: "BaseURL",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <base href="/world/">
			       <link href="https://example.com/hello/foo.jpg"
			           rel="preload" as="image">
			       <link href="bar.jpg" rel="preload" as="image">`,
			want: []*preload.Preload{
				pl(`<https://example.com/hello/foo.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/world/bar.jpg>;rel="preload";as="image"`),
			},
		},
		{
			name: "WebFont",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <link href="/fonts/icons.woff2" rel="preload"
			              as="font" type="font/woff2" crossorigin>`,
			want: []*preload.Preload{
				pl(`<https://example.com/fonts/icons.woff2>;rel="preload";as="font";crossorigin;type="font/woff2"`),
			},
		},
		{
			name: "Media",
			url:  "https://example.com/hello/",
			html: `<!doctype html>
			       <link href="small.jpg" rel="preload" as="image"
			             media="(max-width: 600px)">
			       <link href="large.jpg" rel="preload" as="image"
			             media="(min-width: 601px)">`,
			want: []*preload.Preload{
				pl(`<https://example.com/hello/small.jpg>;rel="preload";as="image";media="(max-width: 600px)"`),
				pl(`<https://example.com/hello/large.jpg>;rel="preload";as="image";media="(min-width: 601px)"`),
			},
		},
	}

	extractPreloadTags := htmltask.ExtractPreloadTags()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := makeHTMLResponse(test.url, test.html)
			if err := extractPreloadTags.Run(resp); err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if diff := cmp.Diff(test.want, resp.Preloads); diff != "" {
				t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
