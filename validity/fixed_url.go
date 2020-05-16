// Copyright 2020 Google LLC
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

package validity

import (
	"net/url"

	"github.com/google/webpackager/exchange"
)

// FixedURL uses url as the validity URL for all resources. It is useful,
// for example, if you plan to provide empty validity data (consisting of
// a single byte 0xa0) for all signed exchanges and serve it on a single URL.
//
// url can be relative, in which case the validity URL is resolved relative
// from the physical URL.
func FixedURL(url *url.URL) URLRule {
	return &fixedURL{url}
}

type fixedURL struct {
	url *url.URL
}

func (rule *fixedURL) Apply(physurl *url.URL, resp *exchange.Response, vp exchange.ValidPeriod) (*url.URL, error) {
	return physurl.ResolveReference(rule.url), nil
}
