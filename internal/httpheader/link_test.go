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

package httpheader_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/internal/httpheader"
	"github.com/google/webpackager/internal/urlutil"
)

func TestParseLink(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []*httpheader.Link
	}{
		{
			name: "AbsoluteURL",
			arg:  `<https://example.org/style.css>; rel="preload"; as="style"`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/style.css"),
					Params: map[string]string{"rel": "preload", "as": "style"},
					Header: `<https://example.org/style.css>; rel="preload"; as="style"`,
				},
			},
		},
		{
			name: "RelativeURL",
			arg:  `</style.css>; rel="preload"; as="style"`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("/style.css"),
					Params: map[string]string{"rel": "preload", "as": "style"},
					Header: `</style.css>; rel="preload"; as="style"`,
				},
			},
		},
		{
			name: "Multiple",
			arg: `<https://example.org/style.css>; rel="preload"; as="style",` +
				`<https://example.org/image.png>; rel="preload"; as="image"`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/style.css"),
					Params: map[string]string{"rel": "preload", "as": "style"},
					Header: `<https://example.org/style.css>; rel="preload"; as="style"`,
				},
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/image.png"),
					Params: map[string]string{"rel": "preload", "as": "image"},
					Header: `<https://example.org/image.png>; rel="preload"; as="image"`,
				},
			},
		},
		{
			name: "Compact",
			arg:  `<https://example.org/style.css>;rel=preload;as=style`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/style.css"),
					Params: map[string]string{"rel": "preload", "as": "style"},
					Header: `<https://example.org/style.css>;rel=preload;as=style`,
				},
			},
		},
		{
			name: "NamesLowered",
			arg:  `<https://example.org/style.css>; REL="preload"; As="style"`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/style.css"),
					Params: map[string]string{"rel": "preload", "as": "style"},
					Header: `<https://example.org/style.css>; REL="preload"; As="style"`,
				},
			},
		},
		{
			name: "SpecialChars",
			arg:  `<https://example.org/hello.html>; title="\"Hello, world!\""; rel="next"`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/hello.html"),
					Params: map[string]string{"title": `"Hello, world!"`, "rel": "next"},
					Header: `<https://example.org/hello.html>; title="\"Hello, world!\""; rel="next"`,
				},
			},
		},
		{
			name: "NoValue",
			arg:  `<https://example.org/style.css>; rel="stylesheet"; crossorigin`,
			want: []*httpheader.Link{
				&httpheader.Link{
					URL:    urlutil.MustParse("https://example.org/style.css"),
					Params: map[string]string{"rel": "stylesheet", "crossorigin": ""},
					Header: `<https://example.org/style.css>; rel="stylesheet"; crossorigin`,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := httpheader.ParseLink(test.arg)
			if err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ParseLink(%q) mismatch (-want +got):\n%s", test.arg, diff)
			}
		})
	}
}
