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

// Package httplink defines a representation of Web Linkings.
package httplink

import (
	"fmt"
	"net/url"
	"strings"
)

// Link represents a Web Linking [RFC 8288], aka. the Link HTTP header.
type Link struct {
	URL    *url.URL
	Params LinkParams
}

// NewLink creates and initializes a new Link with the provided URL u and
// the provided rel parameter.
func NewLink(u *url.URL, rel string) *Link {
	p := make(LinkParams, 1)
	p.Set(ParamRel, rel)
	return &Link{u, p}
}

// IsPreload reports whether the Link involves preloading of the resource.
func (l *Link) IsPreload() bool {
	// TODO(yuizumi): Maybe include rel="prefetch" and similar.
	for _, s := range strings.Fields(l.Params.Get(ParamRel)) {
		if strings.EqualFold(s, RelPreload) {
			return true
		}
	}
	return false
}

// String serializes the Link as it appears in the HTTP header.
func (l *Link) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "<%s>", l.URL)
	l.Params.write(&sb)
	return sb.String()
}

// GoString implements the GoStringer interface.
func (l *Link) GoString() string {
	return fmt.Sprintf("&httplink.Link{URL:&%#v, Params:%#v}", *l.URL, l.Params)
}

// Equal reports whether l and m have the same URL and parameters. The URLs
// are compared by strings.
func (l *Link) Equal(m *Link) bool {
	return l.URL.String() == m.URL.String() && l.Params.Equal(m.Params)
}
