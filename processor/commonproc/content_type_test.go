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
	"testing"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor/commonproc"
)

func TestContentTypeProcessor(t *testing.T) {
	tests := []struct {
		name string
		url  string
		resp string
		want string
	}{
		{
			name: "ContentTypeAdded",
			url:  "https://example.org/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=604800\r\n",
				"Content-Length: 35\r\n",
				"\r\n",
				"<!doctype html><p>Hello, world!</p>",
			),
			want: "text/html; charset=utf-8",
		},
		{
			name: "ContentTypeNotAltered",
			url:  "https://example.org/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=604800\r\n",
				"Content-Length: 35\r\n",
				"Content-Type: text/plain; charset=us-ascii\r\n",
				"\r\n",
				"<!doctype html><p>Hello, world!</p>",
			),
			want: "text/plain; charset=us-ascii",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeResponse(test.url, test.resp)

			if err := commonproc.ContentTypeProcessor.Process(resp); err != nil {
				t.Errorf("got error(%q), want success", err)
			}
			if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf(`resp.Header.Get("X-Content-Type-Options") = %q, want %q`, got, "nosniff")
			}
			if got := resp.Header.Get("Content-Type"); got != test.want {
				t.Errorf(`resp.Header.Get("Content-Type") = %q, want %q`, got, test.want)
			}
		})
	}
}
