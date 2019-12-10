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
	"net/http"
	"net/url"
	"time"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/resource"
	multierror "github.com/hashicorp/go-multierror"
)

var (
	isRedirectCode = map[int]bool{
		http.StatusMovedPermanently:  true, // 301
		http.StatusFound:             true, // 302
		http.StatusSeeOther:          true, // 303
		http.StatusTemporaryRedirect: true, // 307
		http.StatusPermanentRedirect: true, // 308
	}

	errReferenceLoop = errors.New("detected cyclic reference")
)

type packagerTaskRunner struct {
	*Packager

	period exchange.ValidPeriod
	errs   *multierror.Error
	active map[string]bool // Keyed by URLs.
}

func newTaskRunner(p *Packager, vp exchange.ValidPeriod) *packagerTaskRunner {
	return &packagerTaskRunner{
		p,
		vp,
		new(multierror.Error),
		make(map[string]bool),
	}
}

func (runner *packagerTaskRunner) err() error {
	return runner.errs.ErrorOrNil()
}

func (runner *packagerTaskRunner) run(parent *packagerTask, req *http.Request, r *resource.Resource) {
	url := r.RequestURL.String()
	var err error

	if runner.active[url] {
		err = errReferenceLoop
	} else {
		log.Printf("processing %v ...", url)
		runner.active[url] = true
		err = (&packagerTask{runner, parent, req, r}).run()
		delete(runner.active, url)
	}

	if err != nil {
		err = WrapError(err, r.RequestURL)
		runner.errs = multierror.Append(runner.errs, err)
		log.Print(err)
	}
}

type packagerTask struct {
	*packagerTaskRunner

	parent   *packagerTask
	request  *http.Request
	resource *resource.Resource
}

func (task *packagerTask) parentRequest() *http.Request {
	if task.parent == nil {
		return nil
	}
	return task.parent.request
}

func (task *packagerTask) date() time.Time {
	return task.period.Date()
}

func (task *packagerTask) run() error {
	r := task.resource

	req := task.request
	if err := task.RequestTweaker.Tweak(req, task.parentRequest()); err != nil {
		return err
	}

	cached, err := task.ResourceCache.Lookup(req)
	if err != nil {
		return err
	}
	if cached != nil {
		if _, err := task.ExchangeFactory.Verify(cached.Exchange, task.date()); err == nil {
			log.Printf("reusing the existing signed exchange for %s", r.RequestURL)
			*r = *cached
			return nil
		} else {
			log.Printf("renewing the signed exchange for %s: %v", r.RequestURL, err)
		}
	}

	rawResp, err := task.FetchClient.Do(req)
	if err != nil {
		return err
	}
	if isRedirectCode[rawResp.StatusCode] {
		dest, err := rawResp.Location()
		if err != nil {
			return err
		}
		r.RedirectURL = dest
		// TODO(yuizumi): Consider allowing redirects for main resources.
		return fmt.Errorf("redirected to %v", dest)
	}

	purl, err := task.getPhysicalURL(r, rawResp)
	if err != nil {
		return err
	}
	r.PhysicalURL = purl

	vurl, err := task.getValidityURL(r, rawResp)
	if err != nil {
		return err
	}
	r.ValidityURL = vurl

	sxg, err := task.createExchange(rawResp)
	if err != nil {
		return err
	}
	if err := r.SetExchange(sxg); err != nil {
		return err
	}

	// TODO(yuizumi): Generate the validity data.

	return task.ResourceCache.Store(r)
}

func (task *packagerTask) getPhysicalURL(r *resource.Resource, resp *http.Response) (*url.URL, error) {
	u := new(url.URL)
	*u = *r.RequestURL
	task.PhysicalURLRule.Rewrite(u, resp.Header)
	return u, nil
}

func (task *packagerTask) getValidityURL(r *resource.Resource, resp *http.Response) (*url.URL, error) {
	return task.ValidityURLRule.Apply(r.PhysicalURL, resp)
}

func (task *packagerTask) createExchange(rawResp *http.Response) (*signedexchange.Exchange, error) {
	sxgResp, err := exchange.NewResponse(rawResp)
	if err != nil {
		return nil, err
	}
	if err := task.Processor.Process(sxgResp); err != nil {
		return nil, err
	}

	for _, p := range sxgResp.Preloads {
		for _, r := range p.Resources() {
			req, err := newGetRequest(r.RequestURL)
			if err != nil {
				return nil, err
			}
			task.packagerTaskRunner.run(task, req, r)
		}
	}

	vu := task.resource.ValidityURL

	sxg, err := task.ExchangeFactory.NewExchange(sxgResp, task.period, vu)
	if err != nil {
		return nil, err
	}
	if _, err := task.ExchangeFactory.Verify(sxg, task.date()); err != nil {
		return nil, err
	}

	return sxg, nil
}
