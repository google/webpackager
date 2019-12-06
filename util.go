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

package webpackager

import (
	"net/http"
	"net/url"
)

func newGetRequest(url *url.URL) (*http.Request, error) {
	// NOTE: We pass the error through to the caller, but this NewRequest
	// call is expected always to succeed: method is http.MethodGet hence
	// always valid; url is an already parsed value.
	return http.NewRequest(http.MethodGet, url.String(), nil)
}
