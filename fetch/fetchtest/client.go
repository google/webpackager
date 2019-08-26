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

package fetchtest

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/google/webpackager/fetch"
)

// FetchClient fetches content from a test server.
type FetchClient struct {
	client   *http.Client
	requests []*http.Request
}

// NewFetchClient creates and initializes a new FetchClient routing any fetch
// requests to the provided test server s.
func NewFetchClient(s *httptest.Server) *FetchClient {
	client := *fetch.DefaultFetchClient

	client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, s.Listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &FetchClient{client: &client}
}

// Requests returns all HTTP requests the FetchClient has received.
func (c *FetchClient) Requests() []*http.Request {
	return c.requests
}

// Do sends an HTTP request to the test server and returns an HTTP response.
func (c *FetchClient) Do(req *http.Request) (*http.Response, error) {
	c.requests = append(c.requests, req)
	return c.client.Do(req)
}
