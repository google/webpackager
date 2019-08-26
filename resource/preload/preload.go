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
	"github.com/google/webpackager/resource"
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
type Preload interface {
	// Header returns the value of Link HTTP header to perform the preload.
	Header() string

	// Resources returns a set of resources referenced by the preload link.
	// It returns multiple resources when the preload offers more than one
	// option, such as images with multi-source ("imagesrcset") or content
	// negotiations ("variants").
	Resources() []*resource.Resource
}
