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
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/internal/httpheader"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/resource"
)

const headerKey = "Link"

// ExtractPreloadHeaders parses the Link HTTP headers of the provided
// exchange.Response and turns the preload links into preload.Preload
// objects. The preload links will be added to the Preloads field and removed
// from the Link HTTP headers. Note they will be eventually added back to
// the Link HTTP headers when the response is turned into a signed exchange.
var ExtractPreloadHeaders processor.Processor = &extractPreloadHeaders{}

type extractPreloadHeaders struct{}

func (*extractPreloadHeaders) Process(resp *exchange.Response) error {
	values, ok := resp.Header[headerKey]
	if !ok {
		return nil
	}
	resp.Header.Del(headerKey)

	for _, value := range values {
		links, err := httpheader.ParseLink(value)
		if err != nil {
			log.Printf("warning: %v -- this header was ignored", err)
			continue
		}
		for _, link := range links {
			if isPreload(link) {
				r := resource.NewResource(resp.Request.URL.ResolveReference(link.URL))
				resp.AddPreload(newRawPreload(r, link.Params))
			} else {
				resp.Header.Add(headerKey, link.Header)
			}
		}
	}

	return nil
}

func isPreload(link *httpheader.Link) bool {
	for _, s := range strings.Fields(link.Params["rel"]) {
		if strings.EqualFold(s, "preload") {
			return true
		}
	}
	return false
}

type rawPreload struct {
	resource *resource.Resource
	header   string
}

func newRawPreload(r *resource.Resource, params map[string]string) *rawPreload {
	// Reconstruct the header value to increase the chance the rawPreload
	// matches equivalent PlainPreload.

	var header strings.Builder

	fmt.Fprintf(&header, "<%s>", r.RequestURL)

	rel := strings.Fields(params["rel"])
	for i := range rel {
		rel[i] = strings.ToLower(rel[i])
	}
	fmt.Fprintf(&header, ";rel=%q", strings.Join(rel, " "))

	keys := make([]string, 0, len(params)-1)
	for key := range params {
		if key == "rel" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := params[key]
		switch key {
		case "as":
			fmt.Fprintf(&header, ";%s=%q", key, strings.ToLower(val))
		case "crossorigin":
			if val == "" || val == "anonymous" {
				fmt.Fprintf(&header, ";%s", key)
			} else {
				fmt.Fprintf(&header, ";%s=%q", key, val)
			}
		default:
			fmt.Fprintf(&header, ";%s=%q", key, val)
		}
	}
	return &rawPreload{r, header.String()}
}

func (rp *rawPreload) String() string {
	return rp.header
}

func (rp *rawPreload) Header() string {
	return rp.header
}

func (rp *rawPreload) Resources() []*resource.Resource {
	return []*resource.Resource{rp.resource}
}
