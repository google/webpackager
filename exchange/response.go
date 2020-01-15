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

package exchange

import (
	"io/ioutil"
	"net/http"

	"github.com/google/webpackager/resource/preload"
)

// These are keys used in ExtraData. They are prefixed with "X-WebPackager"
// to avoid confusion with real HTTP headers.
const (
	// See htmltask.ExtractSubContentTypes.
	SubContentType = "Webpackager-Sub-Content-Type"
)

const linkHeader = "Link"

// Response represents a pre-signed HTTP exchange to make a signed exchange
// from. It is essentially a wrapper around http.Response. Note the request
// is accessible through http.Response.
type Response struct {
	// Response represents the HTTP response this instance is constructed
	// from. Body has been read into Payload and closed. Other fields may
	// also be mutated by processors (Header in particular).
	*http.Response

	// Payload is the content read from Response.Body and possibly modified
	// by processors.
	Payload []byte

	// Preloads represents preload links to add to HTTP headers.
	Preloads []preload.Preload

	// ExtraData contains information extracted from this Response and
	// used inside the program. Processors extract information and append
	// it with an arbitrary key. Subsequent processors, ValidPeriodRules,
	// and ValidityURLRules can reference the information using that key.
	//
	// The signed exchange will not include ExtraData.
	ExtraData http.Header
}

// NewResponse creates and initializes a new Response wrapping resp.
// The new Response takes the ownership of resp: the caller should not use
// resp after this call.
func NewResponse(resp *http.Response) (*Response, error) {
	payload, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	sxgResp := &Response{resp, payload, nil, make(http.Header)}
	return sxgResp, nil
}

// AddPreload adds p to resp.Preloads if p is not already in resp.Preloads,
// and reports whether p was added. It considers Preloads to be equal when
// Header returns an identical string.
func (resp *Response) AddPreload(p preload.Preload) bool {
	for _, q := range resp.Preloads {
		if p.Header() == q.Header() {
			return false
		}
	}
	resp.Preloads = append(resp.Preloads, p)
	return true
}

// GetFullHeader returns a new http.Header containing all header items
// from resp.Header and resp.Preloads. GetFullHeader makes a deep copy of
// resp.Header, thus does not mutate it.
func (resp *Response) GetFullHeader() http.Header {
	header := make(http.Header)

	for key, val := range resp.Header {
		header[key] = make([]string, len(val))
		copy(header[key], val)
	}

	for _, p := range resp.Preloads {
		for _, r := range p.Resources() {
			if r.Integrity != "" {
				header.Add(linkHeader, r.AllowedAltSXGHeader())
			}
		}
		header.Add(linkHeader, p.Header())
	}

	return header
}
