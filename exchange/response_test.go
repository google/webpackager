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

package exchange_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/resource/preload"
	"github.com/google/webpackager/resource/preload/preloadtest"
)

func makeHTTPResponse(payload []byte) *http.Response {
	rw := httptest.NewRecorder()
	rw.WriteHeader(200)
	rw.Write(payload)

	return rw.Result()
}

func TestNewResponse(t *testing.T) {
	html := []byte("<!doctype html><p>Hello, world!</p>")
	rawResp := makeHTTPResponse(html)

	sxgResp, err := exchange.NewResponse(rawResp)
	if err != nil {
		t.Fatalf("got error(%q), want success", err)
	}
	if sxgResp.Response != rawResp {
		t.Errorf("sxgResp.Response = <%p>, want <%p>", sxgResp.Response, rawResp)
	}
	if got := sxgResp.Payload; !bytes.Equal(got, html) {
		t.Errorf("sxgResp.Payload = %q (%d bytes), want %q (%d bytes)", got, len(got), html, len(html))
	}
	if len(sxgResp.Preloads) != 0 {
		t.Errorf("sxgResp.Preloads = %v, want []", sxgResp.Preloads)
	}
}

func TestAddPreload(t *testing.T) {
	pl := preloadtest.NewPreloadForRawURL

	tests := []struct {
		name  string
		pre   []*preload.Preload
		item  *preload.Preload
		post  []*preload.Preload
		added bool
	}{
		{
			name: "Nonempty",
			pre: []*preload.Preload{
				pl("https://example.org/foo.js", preload.AsScript),
				pl("https://example.org/foo.css", preload.AsStyle),
			},
			item: pl("https://example.org/bar.css", preload.AsStyle),
			post: []*preload.Preload{
				pl("https://example.org/foo.js", preload.AsScript),
				pl("https://example.org/foo.css", preload.AsStyle),
				pl("https://example.org/bar.css", preload.AsStyle),
			},
			added: true,
		},
		{
			name: "Empty",
			pre:  nil,
			item: pl("https://example.org/foo.css", preload.AsStyle),
			post: []*preload.Preload{
				pl("https://example.org/foo.css", preload.AsStyle),
			},
			added: true,
		},
		{
			name: "Duplicate",
			pre: []*preload.Preload{
				pl("https://example.org/foo.js", preload.AsScript),
				pl("https://example.org/foo.css", preload.AsStyle),
			},
			item: pl("https://example.org/foo.css", preload.AsStyle),
			post: []*preload.Preload{
				pl("https://example.org/foo.js", preload.AsScript),
				pl("https://example.org/foo.css", preload.AsStyle),
			},
			added: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeEmptyResponse("https://example.org/")
			resp.Preloads = test.pre
			if got := resp.AddPreload(test.item); got != test.added {
				t.Errorf("got %v, want %v", got, test.added)
			}
			if diff := cmp.Diff(test.post, resp.Preloads); diff != "" {
				t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetFullHeader_PreserveNoLinkHeader(t *testing.T) {
	resp := exchangetest.MakeResponse(
		"https://example.org/hello.html",
		fmt.Sprint(
			"HTTP/1.1 200 OK\r\n",
			"Content-Length: 75\r\n",
			"Content-Type: text/html; charset=utf-8\r\n",
			"\r\n",
			`<!doctype html><link rel="stylesheet" href="style.css">`,
			`<p>Hello, world!</p>`,
		),
	)
	p := preloadtest.NewPreloadForRawURL("https://example.org/style.css", preload.AsStyle)
	resp.AddPreload(p)

	_ = resp.GetFullHeader()

	if got, ok := resp.Header["Link"]; ok {
		t.Errorf(`resp.Header["Link"] = %q, want missing`, got)
	}
}
func TestGetFullHeader_PreserveLinkHeader(t *testing.T) {
	resp := exchangetest.MakeResponse(
		"https://example.org/hello.html",
		fmt.Sprint(
			"HTTP/1.1 200 OK\r\n",
			"Content-Length: 35\r\n",
			"Content-Type: text/html; charset=utf-8\r\n",
			"Link: <style.css>; rel=stylesheet\r\n",
			"\r\n",
			`<!doctype html><p>Hello, world!</p>`,
		),
	)
	p := preloadtest.NewPreloadForRawURL("https://example.org/style.css", preload.AsStyle)
	resp.AddPreload(p)

	_ = resp.GetFullHeader()

	want := []string{"<style.css>; rel=stylesheet"}
	if diff := cmp.Diff(want, resp.Header["Link"]); diff != "" {
		t.Errorf("resp.Header[\"Link\"] mismatch (-want +got):\n%s", diff)
	}
}
