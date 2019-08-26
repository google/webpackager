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

// PreloadStylesheets detects <link rel="stylesheet"> in the <head> element
// and adds referenced stylesheets to the Preloads field.
//
// PreloadStylesheets does not include stylesheets that have "alternate" in
// the rel attribute. Those stylesheets are unused in the initial rendering.
// They are not used at all on some unsupported browsers.
func PreloadStylesheets() HTMLTask {
	return &preloadStylesheets{}
}

type preloadStylesheets struct{}

func isStylesheet(n *html.Node) bool {
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
	stylesheet := false
	for _, linkType := range strings.Fields(rel.Val) {
		if strings.EqualFold(linkType, "stylesheet") {
			stylesheet = true
		}
		if strings.EqualFold(linkType, "alternate") {
			return false
		}
	}
	return stylesheet
}

func (*preloadStylesheets) Run(resp *htmldoc.HTMLResponse) error {
	return htmldoc.Traverse(resp.Doc.Head, func(n *html.Node) error {
		if !isStylesheet(n) {
			return nil
		}
		href := resolveURLAttr(htmldoc.FindAttr(n, "href"), resp.Doc)
		if href == nil {
			return nil
		}
		resp.AddPreload(preload.NewPlainPreload(resource.NewResource(href), preload.AsStyle))
		return nil
	})
}
