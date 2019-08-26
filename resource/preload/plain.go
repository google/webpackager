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

package preload

import (
	"fmt"
	"strings"

	"github.com/google/webpackager/resource"
)

// PlainPreload represents an unconditional preload.
type PlainPreload struct {
	// Resource represents the resource to be preloaded.
	*resource.Resource
	// As represents the value of the as attribute.
	As string
	// CrossOrigin represents whether this preload is cross-origin.
	CrossOrigin bool
	// Media represents the media query for this preload directive.
	Media string
	// Type represents the media type of the preloaded resource.
	Type string
}

// NewPlainPreload creates and initializes a new PlainPreload.
func NewPlainPreload(r *resource.Resource, as string) *PlainPreload {
	return &PlainPreload{Resource: r, As: as}
}

// String is an alias of Header, for the convenience with fmt.Print.
func (p *PlainPreload) String() string {
	return p.Header()
}

// Header implements Preload.Header.
func (p *PlainPreload) Header() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, `<%s>`, p.RequestURL)
	fmt.Fprintf(&sb, `;rel="preload"`)

	if p.As != "" {
		fmt.Fprintf(&sb, ";as=%q", p.As)
	}
	if p.CrossOrigin {
		fmt.Fprintf(&sb, ";crossorigin")
	}
	if p.Media != "" {
		fmt.Fprintf(&sb, ";media=%q", p.Media)
	}
	if p.Type != "" {
		fmt.Fprintf(&sb, ";type=%q", p.Type)
	}

	return sb.String()
}

// Resources implements Preload.Resources. It returns a slice with exactly
// one element: p.Resource.
func (p *PlainPreload) Resources() []*resource.Resource {
	return []*resource.Resource{p.Resource}
}
