// Copyright 2021 Google LLC
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

package cache_test

import (
	"net/http"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
)

func makeRequest(rawurl string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func makeResource(rawurl string) *resource.Resource {
	return resource.NewResource(urlutil.MustParse(rawurl))
}
