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

package htmltask

import (
	"log"
	"net/url"

	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"golang.org/x/net/html"
)

func resolveURLAttr(a *html.Attribute, doc *htmldoc.Document) *url.URL {
	if a == nil {
		return nil
	}
	u, err := url.Parse(a.Val)
	if err != nil {
		log.Printf("warning: invalid %v value %q: %v", a.Key, a.Val, err)
		return nil
	}
	return doc.ResolveReference(u)
}
