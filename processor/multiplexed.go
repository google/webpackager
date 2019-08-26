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

package processor

import (
	"log"
	"mime"

	"github.com/google/webpackager/exchange"
)

// MultiplexedProcessor is a map from media types to processors. The map keys
// are normalized to lowercase and do not include parameters (e.g. "text/html",
// not "Text/HTML" or "text/html; charset=utf-8").
type MultiplexedProcessor map[string]Processor

// Process invokes the processor based on Content-Type of the response.
//
// Process normalizes Content-Type to the media type without parameters,
// using mime.ParseMediaType, then looks up the processor in the map using
// the normalized media type as the key. Process does nothing when the map has
// no entry for the given media type.
func (mp MultiplexedProcessor) Process(resp *exchange.Response) error {
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return nil
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil && err != mime.ErrInvalidMediaParameter {
		log.Printf("warning: invalid Content-Type %q: %v", contentType, err)
		return nil
	}
	p := mp[mediaType]
	if p == nil {
		return nil
	}
	return p.Process(resp)
}
