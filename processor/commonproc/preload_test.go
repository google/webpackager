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

package commonproc_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor/commonproc"
)

func TestExtractPreloadHeaders(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		resp         string
		wantPreloads []string
		wantHeader   http.Header
	}{
		{
			name: "AbsoluteURL",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.com/style.css>;rel=\"preload\";as=\"style\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
			),
			wantPreloads: []string{
				`<https://example.com/style.css>;rel="preload";as="style"`,
			},
			wantHeader: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "RelativeURL",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <style.css>;rel=\"preload\";as=\"style\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
			),
			wantPreloads: []string{
				`<https://example.com/style.css>;rel="preload";as="style"`,
			},
			wantHeader: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "NotPreload",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.com/>;rel=\"start\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
			),
			wantPreloads: []string{},
			wantHeader: http.Header{
				"Link": []string{
					`<https://example.com/>;rel="start"`,
				},
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "Multiple_HeaderRepeated",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.com/style.css>;rel=\"preload\";as=\"style\"\r\n",
				"Link: <https://example.com/photo.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/>;rel=\"start\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
				`<!doctype html><link rel="stylesheet" href="style.css">`,
			),
			wantPreloads: []string{
				`<https://example.com/style.css>;rel="preload";as="style"`,
				`<https://example.com/photo.jpg>;rel="preload";as="image"`,
			},
			wantHeader: http.Header{
				"Link": []string{
					`<https://example.com/>;rel="start"`,
				},
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "Multiple_CommaSeparated",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: "+
					"<https://example.com/style.css>;rel=\"preload\";as=\"style\", "+
					"<https://example.com/photo.jpg>;rel=\"preload\";as=\"image\", "+
					"<https://example.com/>;rel=\"start\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
				`<!doctype html><link rel="stylesheet" href="style.css">`,
			),
			wantPreloads: []string{
				`<https://example.com/style.css>;rel="preload";as="style"`,
				`<https://example.com/photo.jpg>;rel="preload";as="image"`,
			},
			wantHeader: http.Header{
				"Link": []string{
					`<https://example.com/>;rel="start"`,
				},
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "CrossOrigin_Anonymous",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.org/world.html>;rel=\"preload\";as=\"document\";"+
					"crossorigin=\"anonymous\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
			),
			wantPreloads: []string{
				`<https://example.org/world.html>;rel="preload";as="document";` +
					`crossorigin`,
			},
			wantHeader: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "CrossOrigin_UseCredentials",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.org/world.html>;rel=\"preload\";as=\"document\";"+
					"crossorigin=\"use-credentials\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
			),
			wantPreloads: []string{
				`<https://example.org/world.html>;rel="preload";as="document";` +
					`crossorigin="use-credentials"`,
			},
			wantHeader: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "Normalization",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.com/font.woff2>; Type=\"font/woff2\""+
					"; Rel=\" Preload \"; CrossOrigin; As=Font\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
			),
			wantPreloads: []string{
				`<https://example.com/font.woff2>;rel="preload";as="font";` +
					`crossorigin;type="font/woff2"`,
			},
			wantHeader: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeResponse(test.url, test.resp)

			if err := commonproc.ExtractPreloadHeaders.Process(resp); err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			gotPreloads := make([]string, len(resp.Preloads))
			for i, p := range resp.Preloads {
				gotPreloads[i] = p.Header()
			}
			if diff := cmp.Diff(test.wantPreloads, gotPreloads); diff != "" {
				t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantHeader, resp.Header); diff != "" {
				t.Errorf("resp.Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
