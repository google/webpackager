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

func TestRemoveUncachedHeaders(t *testing.T) {
	respText := fmt.Sprint(
		"HTTP/1.1 200 OK\r\n",
		"Cache-Control: public, max-age=604800\r\n",
		"Connection: keep-alive\r\n",
		"Content-Length: 35\r\n",
		"Content-Type: text/html;charset=utf-8\r\n",
		"Keep-Alive: timeout=5, max-1000\r\n",
		"Set-Cookie: id=0123456789abcdef\r\n",
		"\r\n",
		"<!doctype html><p>Hello, world!</p>",
	)
	resp := exchangetest.MakeResponse(
		"https://example.org/hello.html", respText)

	if err := commonproc.RemoveUncachedHeaders.Process(resp); err != nil {
		t.Errorf("got error(%q), want success", err)
	}

	want := http.Header{
		"Cache-Control":  []string{"public, max-age=604800"},
		"Content-Length": []string{"35"},
		"Content-Type":   []string{"text/html;charset=utf-8"},
	}

	if diff := cmp.Diff(want, resp.Header); diff != "" {
		t.Errorf("resp.Header mismatch (-want +got):\n%s", diff)
	}
}
