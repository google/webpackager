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
	"reflect"
	"testing"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/preload"
)

func makeHTTPResponse(payload []byte) *http.Response {
	rw := httptest.NewRecorder()
	rw.WriteHeader(200)
	rw.Write(payload)

	return rw.Result()
}

func plain(rawurl, as string) preload.Preload {
	r := resource.NewResource(urlutil.MustParse(rawurl))
	return preload.NewPlainPreload(r, as)
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
	tests := []struct {
		name  string
		pre   []preload.Preload
		item  preload.Preload
		post  []preload.Preload
		added bool
	}{
		{
			name: "Nonempty",
			pre: []preload.Preload{
				plain("https://example.org/foo.js", preload.AsScript),
				plain("https://example.org/foo.css", preload.AsStyle),
			},
			item: plain("https://example.org/bar.css", preload.AsStyle),
			post: []preload.Preload{
				plain("https://example.org/foo.js", preload.AsScript),
				plain("https://example.org/foo.css", preload.AsStyle),
				plain("https://example.org/bar.css", preload.AsStyle),
			},
			added: true,
		},
		{
			name: "Empty",
			pre:  nil,
			item: plain("https://example.org/foo.css", preload.AsStyle),
			post: []preload.Preload{
				plain("https://example.org/foo.css", preload.AsStyle),
			},
			added: true,
		},
		{
			name: "Duplicate",
			pre: []preload.Preload{
				plain("https://example.org/foo.js", preload.AsScript),
				plain("https://example.org/foo.css", preload.AsStyle),
			},
			item: plain("https://example.org/foo.css", preload.AsStyle),
			post: []preload.Preload{
				plain("https://example.org/foo.js", preload.AsScript),
				plain("https://example.org/foo.css", preload.AsStyle),
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
			if !reflect.DeepEqual(resp.Preloads, test.post) {
				t.Errorf("got %v, want %v", resp.Preloads, test.post)
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
	resp.AddPreload(
		plain("https://example.org/style.css", preload.AsStyle))

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
	resp.AddPreload(
		plain("https://example.org/style.css", preload.AsStyle))

	_ = resp.GetFullHeader()

	want := []string{"<style.css>; rel=stylesheet"}
	if got := resp.Header["Link"]; !reflect.DeepEqual(got, want) {
		t.Errorf(`resp.Header["Link"] = %q, want %q`, got, want)
	}
}
