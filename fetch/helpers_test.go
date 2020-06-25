// Copyright 2020 Google LLC
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

package fetch_test

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

type stubFetcher struct{}

func (s *stubFetcher) Do(req *http.Request) (*http.Response, error) {
	respText := fmt.Sprint(
		"HTTP/1.1 200 OK\r\n",
		"Content-Length: 27\r\n",
		"Content-Type: text/html; charset=utf-8\r\n",
		"\r\n",
		`<!doctype html><p>hello</p>`,
	)
	r := bufio.NewReader(strings.NewReader(respText))
	return http.ReadResponse(r, req)
}
