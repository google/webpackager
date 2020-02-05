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
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"github.com/google/webpackager/resource/httplink"
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

func (*extractPreloadTags) Run(resp *htmldoc.HTMLResponse) error {
	return htmldoc.Traverse(resp.Doc.Root, func(n *html.Node) error {
		if n.Type != html.ElementNode || n.DataAtom != atom.Link {
			return nil
		}
		href := resolveURLAttr(htmldoc.FindAttr(n, "href"), resp.Doc)
		if href == nil {
			return nil
		}

		link := httplink.NewLink(href, "")

		if a := htmldoc.FindAttr(n, "rel"); a != nil {
			link.Params.Set(httplink.ParamRel, a.Val)
		}
		if a := htmldoc.FindAttr(n, "as"); a != nil {
			link.Params.Set(httplink.ParamAs, a.Val)
		}
		if a := htmldoc.FindAttr(n, "crossorigin"); a != nil {
			link.Params.Set(httplink.ParamCrossOrigin, a.Val)
		}
		if a := htmldoc.FindAttr(n, "media"); a != nil {
			link.Params.Set(httplink.ParamMedia, a.Val)
		}
		if a := htmldoc.FindAttr(n, "type"); a != nil {
			link.Params.Set(httplink.ParamType, a.Val)
		}

		if link.IsPreload() {
			resp.AddPreload(preload.NewPreloadForLink(link))
		}

		return nil
	})
}
