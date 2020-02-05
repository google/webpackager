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

package htmldoc_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestFindAttr(t *testing.T) {
	tests := []struct {
		name string
		html string
		key  string
		want *html.Attribute
	}{
		{
			name: "img-src",
			html: `<img src="hello.svg" title="Hello!!">`,
			key:  "src",
			want: &html.Attribute{
				Key: "src",
				Val: "hello.svg",
			},
		},
		{
			name: "img-title",
			html: `<img src="hello.svg" title="Hello!!">`,
			key:  "title",
			want: &html.Attribute{
				Key: "title",
				Val: "Hello!!",
			},
		},
		{
			name: "img-width",
			html: `<img src="hello.svg" title="Hello!!">`,
			key:  "width",
			want: nil,
		},
		{
			name: "IMG-TITLE",
			html: `<IMG Src="hello.svg" Title="Hello!!">`,
			key:  "TITLE",
			// html.Parse() seems to lower the attribute name.
			want: &html.Attribute{
				Key: "title",
				Val: "Hello!!",
			},
		},
		{
			name: "input-required",
			html: `<input name="q" value="hello" required>`,
			key:  "required",
			want: &html.Attribute{
				Key: "required",
			},
		},
		{
			name: "svg-xlink:title",
			html: `<svg xlink:title="Hello!!"></svg>`,
			key:  "xlink:title",
			want: nil, // xlink:title is a foreign key.
		},
		{
			name: "svg-title",
			html: `<svg xlink:title="Hello!!"></svg>`,
			key:  "title",
			want: nil, // xlink:title is a foreign key.
		},
		{
			name: "svg-width",
			html: `<svg width="160" height="160"></svg>`,
			key:  "width",
			want: &html.Attribute{
				Key: "width",
				Val: "160",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(test.html))
			if err != nil {
				t.Fatal(err)
			}
			body := htmldoc.FindNode(doc, atom.Body)
			if body == nil {
				t.Fatal("body not found")
			}
			got := htmldoc.FindAttr(body.FirstChild, test.key)
			if got == test.want {
				return
			}
			if got == nil || test.want == nil || *got != *test.want {
				t.Errorf("got %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestFindNode(t *testing.T) {
	tests := []struct {
		name string
		html string
		tag  atom.Atom
		want string
	}{
		{
			name: "Main",
			html: `<!doctype html>
			       <header><div>header</div></header>
			       <main><div>main</div></main>
			       <footer><div>footer</div></footer>`,
			tag:  atom.Main,
			want: `<main><div>main</div></main>`,
		},
		{
			name: "Div",
			html: `<!doctype html>
			       <header><div>header</div></header>
			       <main><div>main</div></main>
			       <footer><div>footer</div></footer>`,
			tag:  atom.Div,
			want: `<div>header</div>`,
		},
		{
			name: "Span",
			html: `<!doctype html>
			       <header><div>header</div></header>
			       <main><div>main</div></main>
			       <footer><div>footer</div></footer>`,
			tag:  atom.Span,
			want: ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(test.html))
			if err != nil {
				t.Fatal(err)
			}
			node := htmldoc.FindNode(doc, test.tag)
			var got strings.Builder
			if node != nil {
				if err := html.Render(&got, node); err != nil {
					t.Fatal(err)
				}
			}
			if got.String() != test.want {
				t.Errorf("got %q, want %q", got.String(), test.want)
			}
		})
	}
}

func TestTraversal(t *testing.T) {
	errDummy := errors.New("node_test: dummy")

	desc := func(err error) string {
		switch err {
		case nil:
			return "success"
		case errDummy:
			return "errDummy"
		default:
			return fmt.Sprintf("error(%q)", err)
		}
	}

	src := fmt.Sprint(
		`<!doctype html>`,
		`<h1><span>1a;</span><span>1b;</span></h1>`,
		`<h2><span>2a;</span><span>2b;</span></h2>`,
		`<h3><span>3a;</span><span>3b;</span></h3>`,
	)
	out := new(strings.Builder) // Reset at each case.

	tests := []struct {
		name string
		body func(*html.Node) error
		want string
		err  error
	}{
		{
			name: "Base",
			body: func(n *html.Node) error {
				if n.Type == html.TextNode {
					out.WriteString(n.Data)
				}
				return nil // Keep treversal.
			},
			want: "1a;1b;2a;2b;3a;3b;",
			err:  nil,
		},
		{
			name: "Skip",
			body: func(n *html.Node) error {
				if n.Type == html.TextNode {
					out.WriteString(n.Data)
				}
				if n.Type == html.ElementNode && n.DataAtom == atom.H2 {
					return htmldoc.ErrSkip
				}
				return nil // Keep treversal.
			},
			want: "1a;1b;3a;3b;",
			err:  nil,
		},
		{
			name: "Stop",
			body: func(n *html.Node) error {
				if n.Type == html.TextNode {
					out.WriteString(n.Data)
				}
				if n.Type == html.ElementNode && n.DataAtom == atom.H2 {
					return htmldoc.ErrStop
				}
				return nil // Keep treversal.
			},
			want: "1a;1b;",
			err:  nil,
		},
		{
			name: "Fail",
			body: func(n *html.Node) error {
				if n.Type == html.TextNode {
					out.WriteString(n.Data)
				}
				if n.Type == html.ElementNode && n.DataAtom == atom.H2 {
					return errDummy
				}
				return nil // Keep treversal.
			},
			want: "1a;1b;",
			err:  errDummy,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(src))
			if err != nil {
				t.Fatal(err)
			}
			out.Reset()
			if err := htmldoc.Traverse(doc, test.body); err != test.err {
				t.Errorf("got %v, want %v", desc(err), desc(test.err))
			}
			if got := out.String(); got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
