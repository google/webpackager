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

package validity

import (
	"net/url"

	"github.com/google/webpackager/exchange"
)

// URLRule decides the validity URL of a resource.
type URLRule interface {
	// Apply returns the validity URL of a resource. physurl is the physical
	// URL of the resource; resp is the HTTP response; vp is the period the
	// signed exchange will be valid for. The physical URL is typically equal
	// to the request URL but different in some cases; see package urlrewrite
	// for more details.
	//
	// Note implementations can retrieve the request URL via resp.Request.URL.
	Apply(physurl *url.URL, resp *exchange.Response, vp exchange.ValidPeriod) (*url.URL, error)
}

// DefaultURLRule is the default rule used by webpackager.Packager.
var DefaultURLRule URLRule = AppendExtDotLastModified(".validity")
