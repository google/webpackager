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

// Package preloadtest provides utilities for preload link testing.
package preloadtest

import (
	"fmt"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource/httplink"
	"github.com/google/webpackager/resource/preload"
)

// NewPreloadForRawURL is like preload.NewPreloadForURL, but takes a URL
// string instead of a url.URL. It panics on error with parsing rawurl for
// ease of use in testing.
func NewPreloadForRawURL(rawurl, as string) *preload.Preload {
	return preload.NewPreloadForURL(urlutil.MustParse(rawurl), as)
}

// NewPreloadForRawLink is like preload.NewPreloadForLink, but takes an
// HTTP header value instead of an httplink.Link. rawLink must contain
// exactly one valid Web Linking; otherwise NewPreloadForRawLink panics.
func NewPreloadForRawLink(rawLink string) *preload.Preload {
	links, err := httplink.Parse(rawLink)
	if err != nil {
		panic(err)
	}
	if len(links) != 1 {
		panic(fmt.Sprintf("includes %v links: %q", len(links), rawLink))
	}
	return preload.NewPreloadForLink(links[0])
}
