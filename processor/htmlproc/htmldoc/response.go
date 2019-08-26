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

package htmldoc

import (
	"github.com/google/webpackager/exchange"
)

// HTMLResponse is an extension of exchange.Response for HTML responses.
type HTMLResponse struct {
	*exchange.Response
	Doc *Document
}

// NewHTMLResponse creates and initializes a new HTMLResponse.
func NewHTMLResponse(resp *exchange.Response) (*HTMLResponse, error) {
	doc, err := NewDocument(resp.Payload, resp.Request.URL)
	if err != nil {
		return nil, err
	}
	return &HTMLResponse{resp, doc}, nil
}
