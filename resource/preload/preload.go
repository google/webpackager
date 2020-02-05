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

// Package preload defines representations of preload links.
package preload

import (
	"net/url"

	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/httplink"
)

// Values for the "as" attribute of preload links.
const (
	AsAudio    = "audio"
	AsDocument = "document"
	AsEmbed    = "embed"
	AsFetch    = "fetch"
	AsFont     = "font"
	AsImage    = "image"
	AsObject   = "object"
	AsScript   = "script"
	AsStyle    = "style"
	AsTrack    = "track"
	AsWorker   = "worker"
	AsVideo    = "video"
)

// Preload represents a preload link.
type Preload struct {
	*httplink.Link

	// Resources contains a set of resources referenced by the preload link.
	// It typically consists just of one resource, but contains multiple
	// resources when the preload offers more than one option, such as images
	// with multi-source ("imagesrcset") or content negotiations ("variants").
	Resources []*resource.Resource
}

// NewPreloadForURL creates and initializes a new Preload to preload u.
// The new Preload is populated with a new single Resource requesting to u.
// Note it implies u should be absolute.
//
// as specifies the "as" parameter value. If it is empty, the parameter is
// kept unset.
func NewPreloadForURL(u *url.URL, as string) *Preload {
	link := httplink.NewLink(u, httplink.RelPreload)
	if as != "" {
		link.Params.Set(httplink.ParamAs, as)
	}
	return &Preload{link, []*resource.Resource{resource.NewResource(u)}}
}

// NewPreloadForLink creates and initializes a new Preload to perform
// the preloading as specified by link. The new Preload is populated with
// a new single Resource requesting to link.URL. Note it implies link.URL
// should be absolute.
//
// NewPreloadForLink assumes link.IsPreload() to be true.
func NewPreloadForLink(link *httplink.Link) *Preload {
	r := resource.NewResource(link.URL)
	return &Preload{link, []*resource.Resource{r}}
}

// NewPreloadForResource creates and initializes a new Preload to preload
// a single Resource.
//
// as specifies the "as" parameter value. If it is empty, the parameter is
// kept unset.
func NewPreloadForResource(r *resource.Resource, as string) *Preload {
	link := httplink.NewLink(r.RequestURL, httplink.RelPreload)
	if as != "" {
		link.Params.Set(httplink.ParamAs, as)
	}
	return &Preload{link, []*resource.Resource{r}}
}
