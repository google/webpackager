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

package webpackager

import (
	"net/http"
	"time"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/defaultproc"
	"github.com/google/webpackager/resource/cache"
	"github.com/google/webpackager/urlrewrite"
	"github.com/google/webpackager/validity"
)

const validityExt = ".validity"

// Config defines injection points to Packager.
type Config struct {
	// RequestHeader specifies HTTP headers added to every request issued
	// from Packager.
	//
	// Note empty RequestHeader does not imply requests sent without any
	// header fields: http.NewRequest sets a few headers automatically
	// (such as "User-Agent"). RequestHeader takes the precedence over
	// those automatic headers; setting their value to nil suppresses the
	// automatic headers in particular.
	//
	// nil implies empty Header.
	RequestHeader http.Header

	// FetchClient specifies how to retrieve the resources which Packager
	// produces the signed exchanges for.
	//
	// nil implies fetch.DefaultFetchClient, which is just an http.Client
	// properly configured.
	FetchClient fetch.FetchClient

	// PhysicalURLRule specifies the rule(s) to simulate the URL rewriting
	// on the server side, such as appending "index.html" to the path when
	// it points to a directory.
	//
	// nil implies urlrewrite.DefaultRules, which contains a reasonable set
	// of rules to simulate static web servers.
	//
	// See package urlrewrite for details.
	PhysicalURLRule urlrewrite.Rule

	// ValidityURLRule specifies the rule to determine the validity URL,
	// where the validity data should be served.
	//
	// nil implies validity.AppendExtDotUnixTime(".validity", time.Now()).
	ValidityURLRule validity.ValidityURLRule

	// Processor specifies the processor(s) applied to each HTTP response
	// before turning it into a signed exchange. The processors make sure
	// the response can be distributed as signed exchanges and optionally
	// adjust the response for optimized page loading.
	//
	// nil implies defaultproc.DefaultProcessor, a composite of relatively
	// conservative processors.
	//
	// See package processor for details.
	Processor processor.Processor

	// ExchangeFactory specifies encoding/signing parameters for producing
	// signed exchanges.
	//
	// ExchangeFactory must be set to non-nil.
	ExchangeFactory *exchange.Factory

	// ResourceCache specifies the cache to store the signed exchanges and
	// the validity data.
	//
	// It is typically initialized with filewrite.NewFileWriteCache so the
	// signed exchanges are saved into files. See the package document for
	// sample usage.
	//
	// nil implies cache.NewOnMemoryCache(). It is not likely useful:
	// the process would produce signed exchanges and store them in memory,
	// then throw them away at the termination.
	ResourceCache cache.ResourceCache
}

func (cfg *Config) populateDefaults() {
	if cfg.ExchangeFactory == nil {
		panic("ExchangeFactory can't be nil")
	}
	if cfg.RequestHeader == nil {
		cfg.RequestHeader = make(http.Header)
	}
	if cfg.FetchClient == nil {
		cfg.FetchClient = fetch.DefaultFetchClient
	}
	if cfg.PhysicalURLRule == nil {
		cfg.PhysicalURLRule = urlrewrite.DefaultRules
	}
	if cfg.ValidityURLRule == nil {
		cfg.ValidityURLRule = validity.AppendExtDotUnixTime(validityExt, time.Now())
	}
	if cfg.Processor == nil {
		cfg.Processor = defaultproc.DefaultProcessor
	}
	if cfg.ResourceCache == nil {
		cfg.ResourceCache = cache.NewOnMemoryCache()
	}
}
