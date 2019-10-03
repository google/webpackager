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
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/resource"
)

var (
	isRedirectCode = map[int]bool{
		http.StatusMovedPermanently:  true, // 301
		http.StatusFound:             true, // 302
		http.StatusSeeOther:          true, // 303
		http.StatusTemporaryRedirect: true, // 307
		http.StatusPermanentRedirect: true, // 308
	}
)

type packagerTask struct {
	*Packager
	resource *resource.Resource
	period   *exchange.ValidPeriod
	Done     bool
}

func newPackagerTask(pkg *Packager, r *resource.Resource, vp *exchange.ValidPeriod) *packagerTask {
	return &packagerTask{pkg, r, vp, false}
}

func (task *packagerTask) date() time.Time {
	return task.period.Date()
}

func (task *packagerTask) run() error {
	defer (func() { task.Done = true })()
	r := task.resource

	req := task.makeRequest(r.RequestURL)

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

func (task *packagerTask) makeRequest(url *url.URL) *http.Request {
	// NewRequest is expected always successful: method is http.MethodGet
	// thus always valid; url is already parsed value.
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		panic(err)
	}
	for k, v := range task.RequestHeader {
		req.Header[k] = v
	}
	return req
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
			task.Packager.run(r, task.period)
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
