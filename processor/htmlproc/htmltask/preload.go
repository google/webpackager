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
	"strings"

	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/preload"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ExtractPreloadTags detects <link rel="preload"> in the <head> element and
// adds them to the Preloads field.
func ExtractPreloadTags() HTMLTask {
	return &extractPreloadTags{}
}

type extractPreloadTags struct{}

func isPreload(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if n.DataAtom != atom.Link {
		return false
	}
	rel := htmldoc.FindAttr(n, "rel")
	if rel == nil {
		return false
	}
	for _, r := range strings.Fields(rel.Val) {
		if strings.EqualFold(r, "preload") {
			return true
		}
	}
	return false
}

func (*extractPreloadTags) Run(resp *htmldoc.HTMLResponse) error {
	return htmldoc.Traverse(resp.Doc.Root, func(n *html.Node) error {
		if !isPreload(n) {
			return nil
		}
		href := resolveURLAttr(htmldoc.FindAttr(n, "href"), resp.Doc)
		if href == nil {
			return nil
		}
		resp.AddPreload(&preload.PlainPreload{
			Resource:    resource.NewResource(href),
			As:          htmldoc.GetAttr(n, "as"),
			CrossOrigin: (htmldoc.FindAttr(n, "crossorigin") != nil),
			Media:       htmldoc.GetAttr(n, "media"),
			Type:        htmldoc.GetAttr(n, "type"),
		})
		return nil
	})
}
