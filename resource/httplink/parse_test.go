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

package httplink_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource/httplink"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []*httplink.Link
	}{
		{
			name: "AbsoluteURL",
			arg:  `<https://example.org/style.css>; rel="preload"; as="style"`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/style.css"),
					Params: httplink.LinkParams{
						"rel": "preload", "as": "style",
					},
				},
			},
		},
		{
			name: "RelativeURL",
			arg:  `</style.css>; rel="preload"; as="style"`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("/style.css"),
					Params: httplink.LinkParams{
						"rel": "preload", "as": "style",
					},
				},
			},
		},
		{
			name: "Multiple",
			arg: `<https://example.org/style.css>; rel="preload"; as="style",` +
				`<https://example.org/image.png>; rel="preload"; as="image"`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/style.css"),
					Params: httplink.LinkParams{
						"rel": "preload", "as": "style",
					},
				},
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/image.png"),
					Params: httplink.LinkParams{
						"rel": "preload", "as": "image",
					},
				},
			},
		},
		{
			name: "Compact",
			arg:  `<https://example.org/style.css>;rel=preload;as=style`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/style.css"),
					Params: httplink.LinkParams{
						"rel": "preload", "as": "style",
					},
				},
			},
		},
		{
			name: "Normalize",
			arg:  `<https://example.org/style.css>; REL="PRELOAD"; As="Style"`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/style.css"),
					Params: httplink.LinkParams{
						"rel": "preload", "as": "style",
					},
				},
			},
		},
		{
			name: "SpecialChars",
			arg:  `<https://example.org/hello.html>; title="\"Hello, world!\""; rel="next"`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/hello.html"),
					Params: httplink.LinkParams{
						"title": `"Hello, world!"`, "rel": "next",
					},
				},
			},
		},
		{
			name: "NoValue",
			arg:  `<https://example.org/style.css>; rel="stylesheet"; crossorigin`,
			want: []*httplink.Link{
				&httplink.Link{
					URL: urlutil.MustParse("https://example.org/style.css"),
					Params: httplink.LinkParams{
						"rel": "stylesheet", "crossorigin": "anonymous",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := httplink.Parse(test.arg)
			if err != nil {
				t.Fatalf("Parse(%q) returns error(%q), want success", test.arg, err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", test.arg, diff)
			}
		})
	}
}
