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

package webpackager_test

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/webpackager"
	"github.com/google/webpackager/fetch/fetchtest"
)

func stubErrorHandler(status int) http.Handler {
	html := fmt.Sprintf("<!doctype html><p>HTTP error %d</p>", status)
	return stubHandler(status, html, "text/html; charset=utf-8")
}

func stubHTMLHandler(html string) http.Handler {
	return stubHandler(http.StatusOK, html, "text/html; charset=utf-8")
}

func stubTextHandler(text, ctype string) http.Handler {
	return stubHandler(http.StatusOK, text, ctype)
}

func stubHandler(status int, text, ctype string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Length", fmt.Sprint(len(text)))
			w.Header().Set("Content-Type", ctype)
			w.Header().Set("Cache-Control", "public, max-age=1209600")
			w.Header().Set("Date", "Mon, 13 May 2019 10:15:00 GMT")
			w.Header().Set("Expires", "Mon, 27 May 2019 10:15:00 GMT")
			w.WriteHeader(status)
			w.Write([]byte(text))
		},
	)
}

func verifyExchange(t *testing.T, pkg *webpackager.Packager, url, link string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Error(err)
		return
	}
	r, err := pkg.ResourceCache.Lookup(req)
	if err != nil {
		t.Errorf("Lookup(%q) = error(%q), want success", url, err)
		return
	}
	if r == nil {
		t.Errorf("Lookup(%q) = <nil>, want non-nil", url)
		return
	}
	if r.Exchange == nil {
		t.Errorf("Lookup(%q).Exchange = <nil>, want non-nil", url)
		return
	}

	if _, err := pkg.ExchangeFactory.Verify(r.Exchange); err != nil {
		t.Errorf("Verify(sxg[%q]) = error(%q), want success", url, err)
	}
	if got := strings.Join(r.Exchange.ResponseHeaders["Link"], ","); got != link {
		t.Errorf(`sxg[%q].ResponseHeaders.Get("Link") = %#q, want %#q`, url, got, link)
	}
}

func verifyRequests(t *testing.T, pkg *webpackager.Packager, want []string) {
	t.Helper()

	reqs := pkg.FetchClient.(*fetchtest.FetchClient).Requests()

	got := make([]string, len(reqs))
	for i, req := range reqs {
		got[i] = req.URL.String()
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("[Received request URLs] = %q, want %q", got, want)
	}
}
