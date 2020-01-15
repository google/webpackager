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
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	mimeTypeCSS    = "text/style"
	mimeTypeJS     = "application/javascript"
	mimeTypeMathML = "application/mathml+xml"
	mimeTypeSVG    = "image/svg+xml"
)

// ExtractSubContentTypes detects internal subcontents in HTML docuemnt,
// such as CSS (<style>) and JavaScript (<script> without src attribute),
// and adds their MIME types (e.g. "text/style", "application/javascript")
// to ExtraData, using exchange.SubContentType as the key.
func ExtractSubContentTypes() HTMLTask {
	return &extractSubContentTypes{}
}

type extractSubContentTypes struct{}

func (*extractSubContentTypes) Run(resp *htmldoc.HTMLResponse) error {
	return htmldoc.Traverse(resp.Doc.Root, func(n *html.Node) error {
		if n.Type != html.ElementNode {
			return nil
		}

		switch n.DataAtom {
		case atom.Math:
			addSubContentType(resp, mimeTypeMathML)

		case atom.Script:
			if htmldoc.FindAttr(n, "src") != nil {
				return nil
			}
			mimeType := htmldoc.GetAttr(n, "type")
			if mimeType == "" {
				mimeType = mimeTypeJS
			}
			addSubContentType(resp, mimeType)

		case atom.Style:
			mimeType := htmldoc.GetAttr(n, "type")
			if mimeType == "" {
				mimeType = mimeTypeCSS
			}
			addSubContentType(resp, mimeType)

		case atom.Svg:
			addSubContentType(resp, mimeTypeSVG)
		}

		return nil
	})
}

func addSubContentType(resp *htmldoc.HTMLResponse, mimeType string) {
	for _, v := range resp.ExtraData[exchange.SubContentType] {
		if v == mimeType {
			return
		}
	}
	resp.ExtraData.Add(exchange.SubContentType, mimeType)
}
