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
	"net/url"
	"time"

	"github.com/google/webpackager/resource"
	"golang.org/x/xerrors"
)

// Packager implements the control flow of Web Packager.
type Packager struct {
	Config
}

// NewPackager creates and initializes a new Packager with the provided
// Config. It panics when config.ExchangeFactory is nil.
func NewPackager(config Config) *Packager {
	config.populateDefaults()
	return &Packager{config}
}

// Run runs the process to obtain the signed exchange for url: fetches the
// content from the server, processes it, and produces the signed exchange
// from it. Run also takes care of subresources (external resources
// referenced from the fetched content), such as stylesheets, provided they
// are good for preloading.
//
// The process stops when it encounters any error with processing the main
// resource (specified by url), but keeps running and produces the signed
// exchange for the main resource if it just fails with the subresources.
// In either case, Run returns a multierror.Error (hashicorp/go-multierror),
// which consists of webpackager.Errors.
//
// Run does not run the process when ResourceCache already has an entry for
// url.
func (pkg *Packager) Run(url *url.URL, sxgDate time.Time) error {
	req, err := newGetRequest(url)
	if err != nil {
		return err
	}
	return pkg.RunForRequest(req, sxgDate)
}

// RunForRequest is like Run, but takes an http.Request instead of a URL
// thus provides more flexibility to the caller.
//
// RunForRequest uses req directly: RequestTweaker mutates req; FetchClient
// sends req to retrieve the HTTP response.
func (pkg *Packager) RunForRequest(req *http.Request, sxgDate time.Time) error {
	runner, err := newTaskRunner(pkg, sxgDate)
	if err != nil {
		return xerrors.Errorf("packaging: %w", err)
	}
	runner.run(nil, req, resource.NewResource(req.URL))
	return runner.err()
}
