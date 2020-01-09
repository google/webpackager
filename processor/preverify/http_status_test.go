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

package preverify_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/preverify"
)

func TestHTTPStatusCode_Success(t *testing.T) {
	// These tests include http.StatusNoContent (204) solely for testing
	// purpose. It does not mean to recommend producing signed exchanges
	// from a 204 response, in particular.
	tests := []struct {
		name string
		url  string
		proc processor.Processor
		resp string
	}{
		{
			name: "OK",
			url:  "https://example.org/hello.html",
			proc: preverify.HTTPStatusCode(
				http.StatusOK,        // 200
				http.StatusNoContent, // 204
			),
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Length: 35\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>Hello, world!</p>",
			),
		},
		{
			name: "NoContent_Included",
			url:  "https://example.org/hello.html",
			proc: preverify.HTTPStatusCode(
				http.StatusOK,        // 200
				http.StatusNoContent, // 204
			),
			resp: fmt.Sprint(
				"HTTP/1.1 204 No Content\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"\r\n",
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeResponse(test.url, test.resp)
			if err := test.proc.Process(resp); err != nil {
				t.Errorf("got error(%q), want success", err)
			}
		})
	}
}

func TestHTTPStatusCode_Error(t *testing.T) {
	tests := []struct {
		name string
		url  string
		proc processor.Processor
		resp string
	}{
		{
			name: "NoContent_Excluded",
			url:  "https://example.org/hello.html",
			proc: preverify.HTTPStatusCode(
				http.StatusOK,
			),
			resp: fmt.Sprint(
				"HTTP/1.1 204 No Content\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"\r\n",
			),
		},
		{
			name: "NotFound",
			url:  "https://example.org/hello.html",
			proc: preverify.HTTPStatusCode(
				http.StatusOK,        // 200
				http.StatusNoContent, // 204
			),
			resp: fmt.Sprint(
				"HTTP/1.1 404 Not Found\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Length: 35\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>404 Not Found</p>",
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeResponse(test.url, test.resp)
			if err := test.proc.Process(resp); err == nil {
				t.Error("got success, want error")
			}
		})
	}
}
