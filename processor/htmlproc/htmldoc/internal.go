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
	"log"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getBaseURL(head *html.Node) *url.URL {
	base := FindNode(head, atom.Base)
	if base == nil {
		return &url.URL{}
	}
	href := FindAttr(base, "href")
	if href == nil {
		return &url.URL{}
	}
	outcome, err := url.Parse(href.Val)
	if err != nil {
		log.Printf("warning: invalid base url %q: %v", href.Val, err)
		return &url.URL{}
	}
	return outcome
}
