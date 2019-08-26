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
	"errors"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ErrSkip and ErrStop are used by Traverse.
var (
	ErrSkip = errors.New("htmldoc: skip")
	ErrStop = errors.New("htmldoc: stop")
)

// FindAttr returns an Attribute in the given Node with the given key, or nil
// if there is no such attribute. FindAttr inspects only attributes with empty
// Namespace and ignores "foreign attributes."
func FindAttr(n *html.Node, key string) *html.Attribute {
	for _, a := range n.Attr {
		if a.Namespace == "" && strings.EqualFold(a.Key, key) {
			return &a
		}
	}
	return nil
}

// GetAttr returns the value of an attribute in the given Node with the given
// key, or empty string if there is no such attribute. GetAttr considers only
// attributes with empty Namespace and ignores "foreign attributes."
func GetAttr(n *html.Node, key string) string {
	if a := FindAttr(n, key); a != nil {
		return a.Val
	}
	return ""
}

// FindNode locates a descendant of the given Node (including the Node itself)
// which has the given tag. It returns the first descendant when there is more
// than one, and nil when there is none.
func FindNode(n *html.Node, tag atom.Atom) *html.Node {
	if n.DataAtom == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if d := FindNode(c, tag); d != nil {
			return d
		}
	}
	return nil
}

// Traverse performs a pre-order traversal on the parse tree n, calling f on
// each node. f can return ErrSkip to not traverse the subtree of the current
// node and ErrStop to terminate the traversal entirely without error. When
// f returns other non-nil error, Traverse abandons the traversal immediately
// and returns the encountered error.
func Traverse(n *html.Node, f func(*html.Node) error) error {
	if err := traverse(n, f); err != nil && err != ErrStop {
		return err
	}
	return nil
}

func traverse(n *html.Node, f func(*html.Node) error) error {
	if err := f(n); err != nil {
		if err == ErrSkip {
			return nil
		}
		return err
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := traverse(c, f); err != nil {
			return err
		}
	}
	return nil
}
