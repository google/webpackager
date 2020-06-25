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

package certmanager_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"github.com/google/webpackager/internal/certchaintest"
	"github.com/jpillora/backoff"
	"golang.org/x/crypto/ocsp"
)

func jitterlessBackoff() *backoff.Backoff {
	b := certmanager.DefaultBackoff.Copy()
	b.Jitter = false
	return b
}

func stubOCSPHandler(body []byte, header http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Header().Set("Content-Type", "application/ocsp-response")
		if header != nil {
			for key, val := range header {
				w.Header()[key] = val
			}
		}
		w.WriteHeader(http.StatusOK)

		w.Write(body)
	}
}

// TODO(yuizumi): Move these functions into a dedicated package.

func sequencialHandler(handlers ...http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handlers[0].ServeHTTP(w, req)
		handlers = handlers[1:]
	}
}

func makeHTTPClient(s *httptest.Server) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func setRequestCaptor(c *http.Client) *requestCaptor {
	rc := &requestCaptor{nil, c.Transport}
	c.Transport = rc
	return rc
}

type requestCaptor struct {
	Requests []*http.Request
	tripper  http.RoundTripper
}

func (rc *requestCaptor) RoundTrip(req *http.Request) (*http.Response, error) {
	rc.Requests = append(rc.Requests, req)
	return rc.tripper.RoundTrip(req)
}

func TestOCSPClient(t *testing.T) {
	cert0401 := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem")
	ocsp0409 := certchaintest.MustReadOCSPRespFile("../../testdata/ocsp/certmanager_0401_0409.ocsp")

	server := httptest.NewServer(stubOCSPHandler(ocsp0409.Raw, nil))
	defer server.Close()

	client := certmanager.NewOCSPClient(certmanager.OCSPClientConfig{
		HTTPClient:       makeHTTPClient(server),
		RetryPolicy:      jitterlessBackoff(),
		NewFutureEventAt: newDummyFutureEvent,
	})

	reqCaptor := setRequestCaptor(client.HTTPClient)

	now := time.Date(2020, time.April, 9, 15, 0, 0, 0, time.UTC)
	wantNextRun := newDummyFutureEvent(
		time.Date(2020, time.April, 12, 12, 0, 0, 0, time.UTC),
	)
	got, nextRun, err := client.Fetch(cert0401, func() time.Time { return now })

	// Verify client made exactly one request to the desired URL.
	gotURLs := make([]string, len(reqCaptor.Requests))
	for i, req := range reqCaptor.Requests {
		gotURLs[i] = req.URL.String()
	}
	wantURLs := []string{
		"http://ocsp.webpackager-ca.test/MEIwQDA%2BMDwwOjAJBgUrDgMCGgUABBSQ85KB2yrxu3YO8TZixeqpwwUSIAQUPPtYxgu2qJD%2B6Xqj%2BZI57vR%2BcTcCAQs%3D",
	}
	if diff := cmp.Diff(wantURLs, gotURLs); diff != "" {
		t.Errorf("Request URLs mismatch (-want +got):\n%s", diff)
	}

	// nextRun must be valid even if (err != nil).
	if diff := cmp.Diff(wantNextRun, nextRun); diff != "" {
		t.Errorf("nextRun mismatch (-want +got):\n%s", diff)
	}
	if err != nil {
		t.Fatalf("client.Fetch() = error(%q), want success", err)
	}
	if diff := cmp.Diff(ocsp0409, got, ocspComparer); diff != "" {
		t.Errorf("client.Fetch mismatch (-want +got):\n%s", diff)
	}
}

func TestOCSPClientError(t *testing.T) {
	cert0401 := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem")
	ocsp0401 := certchaintest.MustReadOCSPRespFile("../../testdata/ocsp/certmanager_0401_0401.ocsp")

	tests := []struct {
		name    string
		handler http.Handler
	}{
		{
			name:    "HTTPError",
			handler: http.NotFoundHandler(),
		},
		{
			name:    "Invalid",
			handler: stubOCSPHandler([]byte("invalid ocsp"), nil),
		},
		{
			name:    "Expired",
			handler: stubOCSPHandler(ocsp0401.Raw, nil),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(test.handler)
			defer server.Close()

			client := certmanager.NewOCSPClient(certmanager.OCSPClientConfig{
				HTTPClient:       makeHTTPClient(server),
				RetryPolicy:      jitterlessBackoff(),
				NewFutureEventAt: newDummyFutureEvent,
			})

			now := time.Date(2020, time.April, 9, 15, 0, 0, 0, time.UTC)
			wantNextRun := newDummyFutureEvent(
				time.Date(2020, time.April, 9, 15, 0, 1, 0, time.UTC),
			)
			got, nextRun, err := client.Fetch(cert0401, func() time.Time { return now })

			// nextRun must be valid even if (err != nil).
			if diff := cmp.Diff(wantNextRun, nextRun); diff != "" {
				t.Errorf("nextRun mismatch (-want +got):\n%s", diff)
			}
			if err == nil {
				t.Fatalf("client.Fetch() = %#v, want error", got)
			}
		})
	}
}

func TestOCSPClientHTTPCacheExpiry(t *testing.T) {
	cert0401 := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem")

	// cachecontrol internally uses time.Now() to check if the response is
	// cacheable, thus we generate an OCSP response at run time instead of
	// using pre-generated ones.
	thisUpdate := time.Now()
	nextUpdate := thisUpdate.Add(7 * 24 * time.Hour)

	template := ocsp.Response{
		Status:       ocsp.Good,
		SerialNumber: cert0401.Leaf.SerialNumber,
		ThisUpdate:   thisUpdate,
		NextUpdate:   nextUpdate,
	}
	ocspBytes, err := ocsp.CreateResponse(
		cert0401.Issuer,
		cert0401.Issuer,
		template,
		certchaintest.MustReadPrivateKeyFile("../../testdata/CA/inter/key.pem").(*ecdsa.PrivateKey),
	)
	if err != nil {
		t.Fatal(err)
	}
	ocspResp, err := certchain.ParseOCSPResponseForRawChain(ocspBytes, cert0401)
	if err != nil {
		t.Fatal(err)
	}

	const cacheMaxAge = 86400 * time.Second // 24 hours.
	server := httptest.NewServer(stubOCSPHandler(
		ocspBytes,
		http.Header{
			"Cache-Control": []string{"max-age=86400, public, no-transform, must-revalidate"},
			"Date":          []string{thisUpdate.Format(time.RFC850)},
			"Expires":       []string{nextUpdate.Format(time.RFC850)},
			"Last-Modified": []string{thisUpdate.Format(time.RFC850)},
		},
	))
	defer server.Close()

	client := certmanager.NewOCSPClient(certmanager.OCSPClientConfig{
		HTTPClient:       makeHTTPClient(server),
		RetryPolicy:      jitterlessBackoff(),
		NewFutureEventAt: newDummyFutureEvent,
	})

	got, nextRun, err := client.Fetch(cert0401, time.Now)

	// nextRun must be valid even if (err != nil).
	// cachecontrol uses time.Now() to calculate the expiry, thus the exact
	// value is unpredictable. We accept a difference up to 5 seconds.
	wantNextRun := newDummyFutureEvent(time.Now().Add(cacheMaxAge))
	looseComparer := cmp.Comparer(func(x, y time.Time) bool {
		return math.Abs(x.Sub(y).Seconds()) <= 5.0
	})
	if diff := cmp.Diff(wantNextRun, nextRun, looseComparer); diff != "" {
		t.Errorf("nextRun mismatch (-want +got):\n%s", diff)
	}
	if err != nil {
		t.Fatalf("client.Fetch() = error(%q), want success", err)
	}
	if diff := cmp.Diff(ocspResp, got, ocspComparer); diff != "" {
		t.Errorf("client.Fetch mismatch (-want +got):\n%s", diff)
	}
}

func TestOCSPClientRetry(t *testing.T) {
	cert0401 := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem")
	ocsp0409 := certchaintest.MustReadOCSPRespFile("../../testdata/ocsp/certmanager_0401_0409.ocsp")

	handler := sequencialHandler(
		http.NotFoundHandler(),
		http.NotFoundHandler(),
		http.NotFoundHandler(),
		stubOCSPHandler(ocsp0409.Raw, nil),
		http.NotFoundHandler(),
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := certmanager.NewOCSPClient(certmanager.OCSPClientConfig{
		HTTPClient:       makeHTTPClient(server),
		RetryPolicy:      jitterlessBackoff(),
		NewFutureEventAt: newDummyFutureEvent,
	})

	now := time.Date(2020, time.April, 9, 15, 0, 0, 0, time.UTC)
	wantTimeList := []time.Time{
		time.Date(2020, time.April, 9, 15, 0, 1, 0, time.UTC),  // 1st failure
		time.Date(2020, time.April, 9, 15, 0, 3, 0, time.UTC),  // 2nd failure
		time.Date(2020, time.April, 9, 15, 0, 7, 0, time.UTC),  // 3rd failure
		time.Date(2020, time.April, 12, 12, 0, 0, 0, time.UTC), // success
		time.Date(2020, time.April, 12, 12, 0, 1, 0, time.UTC), // 4th failure
	}
	for i, wantTime := range wantTimeList {
		_, nextRun, err := client.Fetch(cert0401, func() time.Time { return now })

		// Verifying err is outside the scope of this test function.
		// We still log it to help diagnose the test failure.
		if err != nil {
			t.Logf("Fetch #%v: %v", i+1, err)
		} else {
			t.Logf("Fetch #%v: success", i+1)
		}

		// nextRun must be valid even if (err != nil).
		nextTime := nextRun.(*dummyFutureEvent).Time
		if !nextTime.Equal(wantTime) {
			t.Errorf("nextRun.Time = %v, want %v", nextTime, wantTime)
		}

		now = nextTime // The next call should wait until nextRun.
	}
}

func TestOCSPClientAllowTestCert(t *testing.T) {
	certSelf := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/self_signed.pem")

	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	client := certmanager.NewOCSPClient(certmanager.OCSPClientConfig{
		AllowTestCert: true,
		// Unrelated, but set just in case.
		HTTPClient:       makeHTTPClient(server),
		RetryPolicy:      jitterlessBackoff(),
		NewFutureEventAt: newDummyFutureEvent,
	})

	now := time.Date(2020, time.April, 9, 15, 0, 0, 0, time.UTC)
	got, nextRun, err := client.Fetch(certSelf, func() time.Time { return now })

	// nextRun must be valid even if (err != nil).
	if _, ok := nextRun.(*futureevent.NeverOccursEvent); !ok {
		t.Errorf("nextRun = %#v, want NeverOccursEvent", nextRun)
	}
	if err != nil {
		t.Fatalf("client.Fetch() = error(%q), want success", err)
	}
	if diff := cmp.Diff(certchain.DummyOCSPResponse, got, ocspComparer); diff != "" {
		t.Errorf("client.Fetch mismatch (-want +got):\n%s", diff)
	}
}
