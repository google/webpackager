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

package preverify

import (
	"net/http"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
)

// HTTPStatusOK ensures the response to have the status code 200 (OK).
var HTTPStatusOK = HTTPStatusCode(http.StatusOK)

// HTTPStatusCode ensures the response to have one of the provided HTTP
// status codes. Its Process method returns an HTTPStatusError on error.
func HTTPStatusCode(expectedCodes ...int) processor.Processor {
	expectedCodeSet := make(map[int]bool, len(expectedCodes))
	for _, code := range expectedCodes {
		expectedCodeSet[code] = true
	}
	return &httpStatusCode{expectedCodeSet}
}

type httpStatusCode struct {
	expected map[int]bool
}

func (h *httpStatusCode) Process(resp *exchange.Response) error {
	if !h.expected[resp.StatusCode] {
		return NewHTTPStatusError(resp.StatusCode)
	}
	return nil
}
