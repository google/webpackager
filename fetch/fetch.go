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

// Package fetch defines interface to retrieve contents to package.
package fetch

import "net/http"

// FetchClient retrieves contents from the server or other data source.
//
// FetchClient should not handle redirects. webpackager.Runner handles redirects
// in its own manner, hence FetchClient should pass any 30x responses through.
//
// An http.Client set up with NeverRedirect, such as DefaultFetchClient, meets
// the contracts and is the most natural choice. Other implementations may
// retrieve contents from other sources such as database or filesystem.
type FetchClient interface {
	// Do handles an HTTP request and returns an HTTP response. It is like
	// http.Client.Do but may not send an HTTP request for real.
	//
	// Do should *not* handle redirects. See the document above.
	Do(req *http.Request) (*http.Response, error)
}

// DefaultFetchClient is a drop-in FetchClient to fetch content via HTTP in
// a usual manner.
var DefaultFetchClient = &http.Client{CheckRedirect: NeverRedirect}

// NeverRedirect instructs http.Client to stop handling the redirect and just
// return the last response instead, when set to the CheckRedirect field.
//
// Technically, NeverRedirect is a function just returning ErrUseLastResponse.
func NeverRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
