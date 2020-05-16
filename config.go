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
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/vprule"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/complexproc"
	"github.com/google/webpackager/resource/cache"
	"github.com/google/webpackager/urlrewrite"
	"github.com/google/webpackager/validity"
)

// Config defines injection points to Packager.
type Config struct {
	// RequestTweaker specifies the mutation applied to every http.Request
	// before it is passed to FetchClient. RequestTweaker is applied
	// both to the http.Request instances passed to the Packager and those
	// generated internally (e.g. for subresources). Note that, however,
	// some RequestTweakers have effect only to subresource requests.
	//
	// nil implies fetch.DefaultRequestTweaker.
	RequestTweaker fetch.RequestTweaker

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
	// nil implies validity.DefaultURLRule, which appends ".validity" plus
	// the last modified time (in UNIX time) to the document URL.
	ValidityURLRule validity.URLRule

	// Processor specifies the processor(s) applied to each HTTP response
	// before turning it into a signed exchange. The processors make sure
	// the response can be distributed as signed exchanges and optionally
	// adjust the response for optimized page loading.
	//
	// nil implies complexproc.DefaultProcessor, a composite of relatively
	// conservative processors.
	//
	// See package processor for details.
	Processor processor.Processor

	// ValidPeriodRule specifies the rule to determine the validity period
	// of signed exchanges.
	//
	// nil implies vprule.DefaultRule.
	ValidPeriodRule vprule.Rule

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
	if cfg.RequestTweaker == nil {
		cfg.RequestTweaker = fetch.DefaultRequestTweaker
	}
	if cfg.FetchClient == nil {
		cfg.FetchClient = fetch.DefaultFetchClient
	}
	if cfg.PhysicalURLRule == nil {
		cfg.PhysicalURLRule = urlrewrite.DefaultRules
	}
	if cfg.ValidityURLRule == nil {
		cfg.ValidityURLRule = validity.DefaultURLRule
	}
	if cfg.Processor == nil {
		cfg.Processor = complexproc.DefaultProcessor
	}
	if cfg.ValidPeriodRule == nil {
		cfg.ValidPeriodRule = vprule.DefaultRule
	}
	if cfg.ResourceCache == nil {
		cfg.ResourceCache = cache.NewOnMemoryCache()
	}
}
