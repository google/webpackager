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

func TestNewLink(t *testing.T) {
	u := urlutil.MustParse("https://example.com/style.css")
	rel := "preload"

	want := &httplink.Link{u, httplink.LinkParams{"rel": rel}}
	got := httplink.NewLink(u, rel)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewLink(%q, %q) mismatch (-want +got):\n%s", u, rel, diff)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		link *httplink.Link
		want string
	}{
		{
			name: "AbsoluteURL",
			link: &httplink.Link{
				URL: urlutil.MustParse("https://example.com/style.css"),
				Params: httplink.LinkParams{
					httplink.ParamRel: "preload",
					httplink.ParamAs:  "style",
				},
			},
			want: `<https://example.com/style.css>;rel="preload";as="style"`,
		},
		{
			name: "RelativeURL",
			link: &httplink.Link{
				URL: urlutil.MustParse("style.css"),
				Params: httplink.LinkParams{
					httplink.ParamRel: "preload",
					httplink.ParamAs:  "style",
				},
			},
			want: `<style.css>;rel="preload";as="style"`,
		},
		{
			name: "CrossOrigin_Anonymous",
			link: &httplink.Link{
				URL: urlutil.MustParse("https://example.org/world.html"),
				Params: httplink.LinkParams{
					httplink.ParamRel:         "preload",
					httplink.ParamCrossOrigin: "anonymous",
					httplink.ParamAs:          "document",
				},
			},
			want: `<https://example.org/world.html>;rel="preload";` +
				`as="document";crossorigin`,
		},
		{
			name: "CrossOrigin_UserCredentials",
			link: &httplink.Link{
				URL: urlutil.MustParse("https://example.org/world.html"),
				Params: httplink.LinkParams{
					httplink.ParamRel:         "preload",
					httplink.ParamCrossOrigin: "user-credentials",
					httplink.ParamAs:          "document",
				},
			},
			want: `<https://example.org/world.html>;rel="preload";` +
				`as="document";crossorigin="user-credentials"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.want, test.link.String()); diff != "" {
				t.Errorf("link.String() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
