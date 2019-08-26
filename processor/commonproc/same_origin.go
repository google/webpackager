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
	"net/url"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/resource/preload"
)

// ApplySameOriginPolicy erases the preload links that may load cross-doamin
// resources.
var ApplySameOriginPolicy processor.Processor = &applySameOriginPolicy{}

type applySameOriginPolicy struct{}

func (*applySameOriginPolicy) Process(resp *exchange.Response) error {
	i := 0

	for _, p := range resp.Preloads {
		if hasSameOrigin(p, resp.Request.URL) {
			resp.Preloads[i] = p
			i++
		}
	}
	resp.Preloads = resp.Preloads[:i]

	return nil
}

func hasSameOrigin(p preload.Preload, referrer *url.URL) bool {
	for _, r := range p.Resources() {
		if !urlutil.HasSameOrigin(r.RequestURL, referrer) {
			return false
		}
	}
	return true
}
