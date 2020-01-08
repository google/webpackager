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
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/preverify"
)

func TestMaxContentLength_Success(t *testing.T) {
	tests := []struct {
		name string
		url  string
		proc processor.Processor
		resp string
	}{
		{
			name: "ClearlySmaller",
			url:  "https://example.org/hello.html",
			proc: preverify.MaxContentLength(48),
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
			name: "ExactlyOnLimit",
			url:  "https://example.org/hello.html",
			proc: preverify.MaxContentLength(48),
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Length: 48\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>abcdefghijklmnopqrstuvwxyz</p>",
			),
		},
		{
			name: "NoContentLengthHeader",
			url:  "https://example.org/hello.html",
			proc: preverify.MaxContentLength(48),
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>Hello, world!</p>",
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

func TestMaxContentLength_Error(t *testing.T) {
	tests := []struct {
		name string
		url  string
		proc processor.Processor
		resp string
	}{
		{
			name: "ClearlyLarger",
			url:  "https://example.org/hello.html",
			proc: preverify.MaxContentLength(48),
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Length: 58\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>abcdefghijklmnopqrstuvwxyz0123456789</p>",
			),
		},
		{
			name: "OneByteLarger",
			url:  "https://example.org/hello.html",
			proc: preverify.MaxContentLength(48),
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Length: 49\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>abcdefghijklmnopqrstuvwxyz!</p>",
			),
		},
		{
			name: "NoContentLengthHeader",
			url:  "https://example.org/hello.html",
			proc: preverify.MaxContentLength(48),
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Cache-Control: public, max-age=1209600\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>abcdefghijklmnopqrstuvwxyz0123456789</p>",
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
