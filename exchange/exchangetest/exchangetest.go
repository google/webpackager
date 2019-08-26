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

// Package exchangetest provides utilities for exchange testing.
package exchangetest

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/google/webpackager/exchange"
)

func httpResponse(url, respText string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(respText)), req)
	if err != nil {
		panic(err)
	}
	return resp
}

// MakeResponse returns a new exchange.Response with a new GET request to url.
// respText is the entire HTTP response.
//
// MakeResponse panics on error for ease of use in testing.
func MakeResponse(url, respText string) *exchange.Response {
	resp, err := exchange.NewResponse(httpResponse(url, respText))
	if err != nil {
		panic(err)
	}
	return resp
}

// MakeEmptyRespnose returns a new exchange.Response with a new GET request to
// url and an empty response with the status code 200 (OK).
//
// MakeEmptyResponse panics on error for ease of use in testing.
func MakeEmptyResponse(url string) *exchange.Response {
	return MakeResponse(url, "HTTP/1.1 200 OK\r\n\r\n")
}
