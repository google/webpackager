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
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/google/webpackager/internal/multierror"
	"github.com/google/webpackager/resource"
)

var (
	errReferenceLoop = errors.New("detected cyclic reference")
)

// Packager implements the control flow of Web Packager.
type Packager struct {
	Config
	tasks map[string]*packagerTask
	errs  *multierror.MultiError
}

// NewPackager creates and initializes a new Packager with the provided
// Config. It panics when config.ExchangeFactory is nil.
func NewPackager(config Config) *Packager {
	config.populateDefaults()
	return &Packager{
		config,
		make(map[string]*packagerTask),
		&multierror.MultiError{},
	}
}

// Run runs the process to obtain the signed exchange for url: fetches the
// content from the server, processes it, and produces the signed exchange
// from it. Run also takes care of subresources (external resources
// referenced from the fetched content), such as stylesheets, provided they
// are good for preloading.
//
// Run does not run the process when ResourceCache already has an entry for
// url.
//
// Errors encountered during the process is accumulated inside the Packager
// and accessible through Err.
func (pkg *Packager) Run(url *url.URL) {
	pkg.run(resource.NewResource(url))
}

// Err returns an error containing all errors encountered in past Run calls.
func (pkg *Packager) Err() error {
	return pkg.errs.Err()
}

func (pkg *Packager) run(r *resource.Resource) {
	if err := pkg.runTask(r); err != nil {
		err = fmt.Errorf("error with processing %v: %v", r.RequestURL, err)
		pkg.errs.Add(err)
		log.Print(err)
	}
}

func (pkg *Packager) runTask(r *resource.Resource) error {
	url := r.RequestURL.String()

	if task := pkg.tasks[url]; task != nil {
		if task.Done {
			*r = *task.resource
			return nil
		} else {
			return errReferenceLoop
		}
	}

	log.Printf("processing %v ...", url)
	pkg.tasks[url] = newPackagerTask(pkg, r)
	return pkg.tasks[url].run()
}
