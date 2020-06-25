// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certmanager

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WICG/webpackage/go/signedexchange/certurl"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"github.com/jpillora/backoff"
	"github.com/pquerna/cachecontrol"
)

// DefaultBackoff is the backoff used by OCSPClient by default.
var DefaultBackoff = backoff.Backoff{
	Factor: 2,
	Jitter: true,
	Min:    time.Second,
	Max:    time.Hour,
}

// DefaultOCSPClient is an OCSPClient with the default configuration.
var DefaultOCSPClient = NewOCSPClient(OCSPClientConfig{})

// OCSPClient represents a client of OCSP over HTTP.
type OCSPClient struct {
	OCSPClientConfig
}

var _ OCSPRespSource = (*OCSPClient)(nil)

// OCSPClientConfig configures OCSPClient.
type OCSPClientConfig struct {
	// HTTPClient is an HTTP client used to send an OCSP request.
	// nil implies http.DefaultClient.
	HTTPClient *http.Client

	// RetryPolicy determines when to make a retry on fetch failure.
	// nil implies DefaultBackoff.
	RetryPolicy *backoff.Backoff

	// AllowTestCert specifies whether to allow certificates without
	// any OCSP URI. If AllowTestCert is set true, OCSPClient returns
	// DummyOCSPResponse for OCSP-less certificates.
	AllowTestCert bool

	// NewFutureEventAt is called by Fetch to create the nextRun return
	// parameter. nil implies futureevent.DefaultFactory.
	NewFutureEventAt futureevent.Factory
}

// NewOCSPClient creates and initializes a new OCSPClient.
func NewOCSPClient(config OCSPClientConfig) *OCSPClient {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	if config.RetryPolicy == nil {
		config.RetryPolicy = DefaultBackoff.Copy()
	}
	if config.NewFutureEventAt == nil {
		config.NewFutureEventAt = futureevent.DefaultFactory
	}
	return &OCSPClient{config}
}

// Fetch sends an OCSP request to the OCSP responder at chain.OCSPServer and
// returns the parsed OCSP response, with a futureevent.Event to notify when
// this OCSPClient expects the next call of Fetch.
//
// now is a function that returns the current time, usually time.Now. Fetch
// calls it after retrieving and parsing the OCSP response and examines the
// validity as of the returned time.
//
// On success, nextRun will be set with the middle point between ThisUpdate
// and NextUpdate of the OCSP response or the cache expiry time of the HTTP
// response, whichever comes earlier. On failure, the next Fetch will be
// scheduled for a retry at the time determined based on c.RetryPolicy.
//
// Keep in mind that a clock skew between the local machine and the OCSP
// server could cause a valid response to be judged invalid. To mitigate it,
// include some tweaks in the now function.
func (c OCSPClient) Fetch(chain *certchain.RawChain, now func() time.Time) (ocspResp *certchain.OCSPResponse, nextRun futureevent.Event, err error) {
	retryWait := c.RetryPolicy.Duration()

	if c.AllowTestCert && chain.OCSPServer == "" {
		// The OCSP response never changes for this RawChain. Augmentor
		// calls Fetch when RawChain gets updated, so NeverOccurs works
		// perfectly.
		return certchain.DummyOCSPResponse, futureevent.NeverOccurs(), nil
	}

	raw, expiry, err := c.fetch(chain)
	if err != nil {
		err = fmt.Errorf("cannot fetch OCSP response: %v", err)
		return nil, c.NewFutureEventAt(now().Add(retryWait)), err
	}

	ocspResp, err = certchain.ParseOCSPResponseForRawChain(raw, chain)
	if err != nil {
		err = fmt.Errorf("cannot parse OCSP response: %v", err)
		return nil, c.NewFutureEventAt(now().Add(retryWait)), err
	}

	t := now()
	if err := ocspResp.VerifyForRawChain(t, chain); err != nil {
		err = fmt.Errorf("invalid OCSP response as of %v: %v", t, err)
		return nil, c.NewFutureEventAt(t.Add(retryWait)), err
	}
	if err := ocspResp.VerifySXGCriteria(); err != nil {
		err = fmt.Errorf("invalid OCSP response as of %v: %v", t, err)
		return nil, c.NewFutureEventAt(t.Add(retryWait)), err
	}

	c.RetryPolicy.Reset()

	var nextRunAt time.Time
	duration := ocspResp.NextUpdate.Sub(ocspResp.ThisUpdate)
	midpoint := ocspResp.ThisUpdate.Add(duration / 2)
	if expiry.After(ocspResp.ThisUpdate) && expiry.Before(midpoint) {
		nextRunAt = expiry
	} else {
		nextRunAt = midpoint
	}

	return ocspResp, c.NewFutureEventAt(nextRunAt), nil
}

func (c OCSPClient) fetch(chain *certchain.RawChain) (body []byte, cacheExpiry time.Time, err error) {
	// cacheExpiry can be zero, which means the response was judged to be
	// uncacheable. It is okay the caller ignores noncacheability.
	//
	// Rationale:
	//
	// RFC 6960 doesn't require any HTTP headers other than Content-Type and
	// Content-Length. In particular, the response doesn't have to include
	// cache directives. Also RFC 6960 doesn't require the client to respect
	// cache directives.
	//
	// RFC 5019 defines cache directives in OCSP responses. It is optional to
	// follow that RFC though.
	//
	// In general we respect cache directives so comply with RFC 5019. However,
	// if the cache expiry time isn't useful for our purpose, we ignore it and
	// assume cache directives to be missing.
	httpReq, err := certurl.CreateOCSPRequest(chain.Certs, true)
	if err != nil {
		return nil, time.Time{}, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		err := fmt.Errorf("error response (%v)", httpResp.Status)
		return nil, time.Time{}, err
	}
	ocspResp, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, time.Time{}, err
	}

	reasons, expiry, err := cachecontrol.CachableResponse(httpReq, httpResp, cachecontrol.Options{})
	if len(reasons) > 0 || err != nil {
		expiry = time.Time{}
	}
	return ocspResp, expiry, nil
}
