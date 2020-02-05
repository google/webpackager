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

package htmltask_test

import (
	"fmt"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
)

// makeHTMLResponse returns a new htmldoc.HTMLResponse with a new GET request.
// html is an HTML document (payload). makeHTMLResponse adds a reasonable set
// of HTTP headers automatically.
//
// makeHTMLResponse panics on error for each of use in testing.
func makeHTMLResponse(url, html string) *htmldoc.HTMLResponse {
	httpResp := fmt.Sprint(
		"HTTP/1.1 200 OK\r\n",
		"Cache-Control: public, max-age=604800\r\n",
		"Content-Length: ", len(html), "\r\n",
		"Content-Type: text/html;charset=utf-8\r\n",
		"\r\n",
		html)
	htmlResp, err := htmldoc.NewHTMLResponse(exchangetest.MakeResponse(url, httpResp))
	if err != nil {
		panic(err)
	}
	return htmlResp
}
