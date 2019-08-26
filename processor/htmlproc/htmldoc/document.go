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
	"bytes"
	"errors"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Document represents an HTML document, holding the parse tree and the related
// information.
type Document struct {
	// Root is the root of the HTML parse tree. Note it is a DocumentNode, not
	// <html> element.
	Root *html.Node

	// Head is the node corresponding to <head> element in the parse tree.
	Head *html.Node

	// Body is the node corresponding to <body> element in the parse tree.
	Body *html.Node

	// URL represents where the document is located.
	URL *url.URL

	// BaseURL represents the base URL of the document. It is usually the same
	// as URL above, but can be altered by <base> element.
	BaseURL *url.URL
}

// NewDocument creates and initializes a new Document from payload and url.
func NewDocument(payload []byte, url *url.URL) (*Document, error) {
	root, err := html.Parse(bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// NOTE(yuizumi): These errors do not happen in practice since html.Parse
	// always adds <head> and <body> when they are missing.
	head := FindNode(root, atom.Head)
	if head == nil {
		return nil, errors.New("missing <head>")
	}
	body := FindNode(root, atom.Body)
	if body == nil {
		return nil, errors.New("missing <body>")
	}

	doc := &Document{root, head, body, url, url.ResolveReference(getBaseURL(head))}
	return doc, nil
}

// ResolveReference resolves a URI reference in Document to an absolute URI.
// The URI reference may be relative or absolute. ResolveReference always
// returns a new URL instance, even if the returned URL is identical to the
// reference.
func (doc *Document) ResolveReference(ref *url.URL) *url.URL {
	return doc.BaseURL.ResolveReference(ref)
}
