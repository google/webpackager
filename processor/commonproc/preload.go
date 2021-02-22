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

package commonproc

import (
	"log"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/resource/httplink"
	"github.com/google/webpackager/resource/preload"
)

const headerKey = "Link"

// ExtractPreloadHeaders parses the Link HTTP headers of the provided
// exchange.Response and turns the preload links into preload.Preload
// objects. The preload links will be added to the Preloads field and removed
// from the Link HTTP headers. Note they will be eventually added back to
// the Link HTTP headers when the response is turned into a signed exchange.
var ExtractPreloadHeaders processor.Processor = &extractPreloadHeaders{}

// KeepNonPreloadLinkHeaders instruct the processor to include preload link
// headers that don't have "preload" as the parameter.
// TODO(banaag): put this in the TOML and propagate the config to the processor.
var keepNonPreloadLinkHeaders = false

// maxNumPreloads is the maximum number of preload links allowed by WebPackager.
// This exists to satisfy: https://github.com/google/webpackager/blob/master/docs/cache_requirements.md.
const maxNumPreloads = 20

type extractPreloadHeaders struct{}

func (*extractPreloadHeaders) Process(resp *exchange.Response) error {
	values, ok := resp.Header[headerKey]
	if !ok {
		return nil
	}
	resp.Header.Del(headerKey)

	numPreloads := 0
	for _, value := range values {
		links, err := httplink.Parse(value)
		if err != nil {
			log.Printf("warning: %v -- this header was ignored", err)
			continue
		}
		for _, link := range links {
			if link.IsPreload() {
				if numPreloads < maxNumPreloads {
					link.URL = resp.Request.URL.ResolveReference(link.URL)
					resp.AddPreload(preload.NewPreloadForLink(link))
					numPreloads++
				}
			} else if keepNonPreloadLinkHeaders {
				resp.Header.Add(headerKey, link.String())
			}
		}
	}

	return nil
}
