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

package htmldoc_test

import (
	"net/url"
	"testing"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
)

func TestDocument(t *testing.T) {
	tests := []struct {
		name string
		html string
		url  *url.URL
		base *url.URL
	}{
		{
			name: "WithoutBase",
			html: `<!doctype html><title>test</title><main>hello</main>`,
			url:  urlutil.MustParse("https://dummy.test/hello/foo.html"),
			base: urlutil.MustParse("https://dummy.test/hello/foo.html"),
		},
		{
			name: "WithBase",
			html: `<!doctype html><base href="../"><title>test</title>` +
				`<main>hello</main>`,
			url:  urlutil.MustParse("https://dummy.test/hello/foo.html"),
			base: urlutil.MustParse("https://dummy.test/"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := htmldoc.NewDocument([]byte(test.html), test.url)
			if err != nil {
				t.Fatalf("got error(%q), want success", err)
			}

			if doc.Root == nil {
				t.Errorf("doc.Root = %v, want non-nil", doc.Root)
			}
			if doc.Head == nil {
				t.Errorf("doc.Head = %v, want non-nil", doc.Head)
			}
			if doc.Body == nil {
				t.Errorf("doc.Body = %v, want non-nil", doc.Body)
			}

			if *doc.URL != *test.url {
				t.Errorf("doc.URL = %q, want %q", doc.URL, test.url)
			}
			if *doc.BaseURL != *test.base {
				t.Errorf("doc.BaseURL = %q, want %q", doc.BaseURL, test.base)
			}
		})
	}
}
