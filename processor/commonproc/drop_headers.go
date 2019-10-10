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
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
)

var (
	uncachedHeaders = []string{
		// Hop-by-hop headers.
		"Connection",
		"Keep-Alive",
		"Proxy-Connection",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",

		// Stateful headers.
		"Authentication-Control",
		"Authentication-Info",
		"Clear-Site-Data",
		"Optional-WWW-Authenticate",
		"Proxy-Authenticate",
		"Proxy-Authentication-Info",
		"Public-Key-Pins",
		"Sec-WebSocket-Accept",
		"Set-Cookie",
		"Set-Cookie2",
		"SetProfile",
		"Strict-Transport-Security",
		"WWW-Authenticate",
	}
)

// RemoveUncachedHeaders removes uncached header fields. Such header fields
// are disallowed in signed exchanges.
// https://tools.ietf.org/html/draft-yasskin-http-origin-signed-responses-07#section-4.1.
var RemoveUncachedHeaders processor.Processor = &removeUncachedHeaders{}

type removeUncachedHeaders struct{}

func (*removeUncachedHeaders) Process(resp *exchange.Response) error {
	// TODO:
	for _, name := range uncachedHeaders {
		resp.Header.Del(name)
	}
	return nil
}
