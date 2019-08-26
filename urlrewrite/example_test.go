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

package urlrewrite_test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/webpackager/urlrewrite"
)

func ExampleCleanPath() {
	examples := []string{
		"https://example.com/foo/../index.html",
		"https://example.com/foo/./bar/",
	}

	for _, example := range examples {
		u, err := url.Parse(example)
		if err != nil {
			panic(err)
		}
		urlrewrite.CleanPath().Rewrite(u, make(http.Header))
		fmt.Println(u)
	}

	// Output:
	// https://example.com/index.html
	// https://example.com/foo/bar/
}
