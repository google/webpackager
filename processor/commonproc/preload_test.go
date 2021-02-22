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
	"github.com/google/webpackager/resource/preload"
	"github.com/google/webpackager/resource/preload/preloadtest"
)

func TestExtractPreloadHeaders(t *testing.T) {
	pl := preloadtest.NewPreloadForRawLink

	tests := []struct {
		name         string
		url          string
		resp         string
		wantPreloads []*preload.Preload
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
			wantPreloads: []*preload.Preload{
				pl(`<https://example.com/style.css>;rel="preload";as="style"`),
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
			wantPreloads: []*preload.Preload{
				pl(`<https://example.com/style.css>;rel="preload";as="style"`),
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
			wantPreloads: nil,
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
			wantPreloads: []*preload.Preload{
				pl(`<https://example.com/style.css>;rel="preload";as="style"`),
				pl(`<https://example.com/photo.jpg>;rel="preload";as="image"`),
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
			wantPreloads: []*preload.Preload{
				pl(`<https://example.com/style.css>;rel="preload";as="style"`),
				pl(`<https://example.com/photo.jpg>;rel="preload";as="image"`),
			},
			wantHeader: http.Header{
				"Link": []string{
					`<https://example.com/>;rel="start"`,
				},
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
		},
		{
			name: "Multiple_MoreThanMaxPreloads",
			url:  "https://example.com/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Link: <https://example.com/photo1.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo2.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo3.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo4.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo5.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo6.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo7.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo8.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo9.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo10.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo11.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo12.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo13.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo14.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo15.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo16.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo17.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo18.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo19.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo20.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/photo21.jpg>;rel=\"preload\";as=\"image\"\r\n",
				"Link: <https://example.com/>;rel=\"start\"\r\n",
				"Content-Type: text/html; charset=utf-8\r\n\r\n",
				`<!doctype html><link rel="stylesheet" href="style.css">`,
			),
			wantPreloads: []*preload.Preload{
				pl(`<https://example.com/photo1.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo2.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo3.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo4.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo5.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo6.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo7.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo8.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo9.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo10.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo11.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo12.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo13.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo14.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo15.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo16.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo17.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo18.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo19.jpg>;rel="preload";as="image"`),
				pl(`<https://example.com/photo20.jpg>;rel="preload";as="image"`),
			},
			wantHeader: http.Header{
				"Link": []string{
					`<https://example.com/>;rel="start"`,
				},
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
			if diff := cmp.Diff(test.wantPreloads, resp.Preloads); diff != "" {
				t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantHeader, resp.Header); diff != "" {
				t.Errorf("resp.Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
