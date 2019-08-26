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
	"testing"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor/preverify"
)

func TestRequireStatusOK_Success(t *testing.T) {
	tests := []struct {
		name string
		url  string
		resp string
	}{
		{
			name: "OK",
			url:  "https://example.org/hello.html",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Length: 35\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>Hello, world!</p>",
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeResponse(test.url, test.resp)
			if err := preverify.RequireStatusOK.Process(resp); err != nil {
				t.Errorf("got error(%q), want success", err)
			}
		})
	}
}

func TestRequireStatusOK_Error(t *testing.T) {
	tests := []struct {
		name string
		url  string
		resp string
	}{
		{
			name: "NonOKStatus",
			url:  "https://example.org/hello.html",
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
			if err := preverify.RequireStatusOK.Process(resp); err == nil {
				t.Error("got success, want error")
			}
		})
	}
}
