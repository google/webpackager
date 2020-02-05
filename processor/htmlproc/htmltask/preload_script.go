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
	"fmt"
	"strings"
	"unicode"

	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"github.com/google/webpackager/resource/preload"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// These elements usually render nothing and never contain <script>.
	skipElements = map[atom.Atom]bool{
		atom.Base:     true,
		atom.Link:     true,
		atom.Meta:     true,
		atom.Noscript: true,
		atom.Style:    true,
		atom.Title:    true,
	}
	// These elements usually render some non-trivial content alone, without
	// any inner element, text, or CSS.
	stopElements = map[atom.Atom]bool{
		atom.Applet:   true,
		atom.Audio:    true,
		atom.Button:   true,
		atom.Embed:    true,
		atom.Iframe:   true,
		atom.Img:      true,
		atom.Input:    true,
		atom.Meter:    true,
		atom.Object:   true,
		atom.Select:   true,
		atom.Textarea: true,
		atom.Video:    true,
	}
)

// InsecurePreloadScripts detects scripts loaded at the top of the document
// in a blocking manner and adds them to the Preloads field. More precisely,
// InsecurePreloadScripts scans over the parsed tree up to the first node that
// renders non-trivial content, such as images and non-whitespace text, as the
// top of the document. Scripts are considered as blocking when loaded without
// async or defer attribute.
//
// SECURITY NOTICE: InsecurePreloadScripts should be used with special care.
// Web Packager produces signed exchanges for all resources to be preloaded.
// Remember that signed exchanges will remain used and distributed until they
// expire, once they are made public: it is possible that they continue to be
// distributed on caches you do not know about, even if they turn out to have
// security issues in JavaScript code. To migitate the risk, it is advised to
// consider making signed exchanges expire shortly (say, within 24 hours), or
// even better not preloading scripts. For example, you can load your scripts
// with defer attribute (e.g. <script defer src="foo.js">) to allow browsers
// to keep rendering the web page without waiting for scripts getting loaded
// and executed, thus eliminate the need for preloading.
func InsecurePreloadScripts() HTMLTask {
	return &preloadScripts{}
}

type preloadScripts struct{}

func (*preloadScripts) Run(resp *htmldoc.HTMLResponse) error {
	return htmldoc.Traverse(resp.Doc.Root, func(n *html.Node) error {
		switch n.Type {
		case html.ElementNode:
			if n.DataAtom == atom.Script {
				handleScript(resp, n)
				return htmldoc.ErrSkip
			}
			if skipElements[n.DataAtom] {
				return htmldoc.ErrSkip
			}
			if stopElements[n.DataAtom] {
				return htmldoc.ErrStop
			}
			return nil // Keep traversing.

		case html.TextNode:
			if strings.IndexFunc(n.Data, isNotSpace) >= 0 {
				return htmldoc.ErrStop
			}
			return htmldoc.ErrSkip

		case html.DocumentNode:
			return nil // Keep traversing.

		case html.ErrorNode, html.CommentNode, html.DoctypeNode:
			return htmldoc.ErrSkip

		default:
			return fmt.Errorf("script: unknown NodeType %v", n.Type)
		}
	})
}

func handleScript(resp *htmldoc.HTMLResponse, n *html.Node) {
	if htmldoc.FindAttr(n, "async") != nil {
		return
	}
	if htmldoc.FindAttr(n, "defer") != nil {
		return
	}
	if u := resolveURLAttr(htmldoc.FindAttr(n, "src"), resp.Doc); u != nil {
		resp.AddPreload(preload.NewPreloadForURL(u, preload.AsScript))
	}
}

func isNotSpace(r rune) bool { return !unicode.IsSpace(r) }
