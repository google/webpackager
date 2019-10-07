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
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WICG/webpackage/go/signedexchange/version"
	"github.com/google/webpackager"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/fetch/fetchtest"
	"github.com/google/webpackager/internal/certutil/certtest"
	"github.com/google/webpackager/internal/urlutil"
)

var (
	date    = time.Date(2019, time.May, 13, 10, 30, 0, 0, time.UTC)
	expires = time.Date(2019, time.May, 20, 10, 30, 0, 0, time.UTC)
)

func makeConfig(server *httptest.Server) webpackager.Config {
	return webpackager.Config{
		FetchClient: fetchtest.NewFetchClient(server),
		ExchangeFactory: &exchange.Factory{
			Version:      version.Version1b3,
			MIRecordSize: 4096,
			CertChain:    certtest.ReadCertChainFile("testdata/certs/test.cbor"),
			CertURL:      urlutil.MustParse("https://example.org/cert.cbor"),
			PrivateKey:   certtest.ReadPrivateKeyFile("testdata/certs/test.key"),
		},
	}
}

func TestSameDomain(t *testing.T) {
	handlers := http.NewServeMux()
	handlers.Handle(
		"example.org/hello.html",
		stubHTMLHandler(`<!doctype html>`+
			`<link href="https://example.org/style.css" rel="stylesheet">`+
			`<p>Hello, world!</p>`),
	)
	handlers.Handle(
		"example.org/style.css",
		stubTextHandler(`body { font-family: sans-serif; }`, "text/css"),
	)
	server := httptest.NewTLSServer(handlers)
	defer server.Close()

	pkg := webpackager.NewPackager(makeConfig(server))
	pkg.Run(urlutil.MustParse("https://example.org/hello.html"),
		exchange.NewValidPeriod(date, expires))

	// style.css is on the same domain thus fetched.
	verifyRequests(t, pkg, []string{
		"https://example.org/hello.html",
		"https://example.org/style.css",
	})
	// Exchanges are generated with preloading.
	verifyExchange(t, pkg, "https://example.org/hello.html", date, fmt.Sprint(
		`<https://example.org/style.css>;rel="allowed-alt-sxg";`+
			`header-integrity="sha256-+Xd20Pyxhd3oSvNo2ucj9gdj7ZkHavIaDGkucYF76J8=",`,
		`<https://example.org/style.css>;rel="preload";as="style"`))
	verifyExchange(t, pkg, "https://example.org/style.css", date, "")
}

func TestCrossDomain(t *testing.T) {
	handlers := http.NewServeMux()
	handlers.Handle(
		"example.org/hello.html",
		stubHTMLHandler(`<!doctype html>`+
			`<link href="https://example.com/style.css" rel="stylesheet">`+
			`<p>Hello, world!</p>`),
	)
	handlers.Handle(
		"example.com/style.css",
		stubTextHandler(`body { font-family: sans-serif; }`, "text/css"),
	)
	server := httptest.NewTLSServer(handlers)
	defer server.Close()

	pkg := webpackager.NewPackager(makeConfig(server))
	pkg.Run(urlutil.MustParse("https://example.org/hello.html"),
		exchange.NewValidPeriod(date, expires))

	// style.css is on a cross origin and not fetched: DefaultProcessor
	// includes RequireSameOrigin.
	verifyRequests(t, pkg, []string{
		"https://example.org/hello.html",
	})
	// An exchange is generated without preloading.
	verifyExchange(t, pkg, "https://example.org/hello.html", date, "")
}

func TestDupResource(t *testing.T) {
	handlers := http.NewServeMux()
	handlers.Handle(
		"example.org/hello.html",
		stubHTMLHandler(`<!doctype html><link href="style.css" rel="stylesheet">`+
			`<p>Hello, world!</p>`),
	)
	handlers.Handle(
		"example.org/quick.html",
		stubHTMLHandler(`<!doctype html><link href="style.css" rel="stylesheet">`+
			`<p>The quick brown fox jumps over the lazy dog.</p>`),
	)
	handlers.Handle(
		"example.org/style.css",
		stubTextHandler(`body { font-family: sans-serif; }`, "text/css"),
	)
	server := httptest.NewTLSServer(handlers)
	defer server.Close()

	pkg := webpackager.NewPackager(makeConfig(server))
	pkg.Run(urlutil.MustParse("https://example.org/hello.html"),
		exchange.NewValidPeriod(date, expires))
	pkg.Run(urlutil.MustParse("https://example.org/quick.html"),
		exchange.NewValidPeriod(date, expires))

	// style.css should be fetched only once.
	verifyRequests(t, pkg, []string{
		"https://example.org/hello.html",
		"https://example.org/style.css",
		"https://example.org/quick.html",
	})
	// Exchanges are generated with preloading.
	verifyExchange(t, pkg, "https://example.org/hello.html", date, fmt.Sprint(
		`<https://example.org/style.css>;rel="allowed-alt-sxg";`+
			`header-integrity="sha256-+Xd20Pyxhd3oSvNo2ucj9gdj7ZkHavIaDGkucYF76J8=",`,
		`<https://example.org/style.css>;rel="preload";as="style"`))
	verifyExchange(t, pkg, "https://example.org/quick.html", date, fmt.Sprint(
		`<https://example.org/style.css>;rel="allowed-alt-sxg";`+
			`header-integrity="sha256-+Xd20Pyxhd3oSvNo2ucj9gdj7ZkHavIaDGkucYF76J8=",`,
		`<https://example.org/style.css>;rel="preload";as="style"`))
	verifyExchange(t, pkg, "https://example.org/style.css", date, "")
}

func TestRequestHeader(t *testing.T) {
	handlers := http.NewServeMux()
	handlers.Handle(
		"example.org/hello.html",
		stubHTMLHandler(`<!doctype html><link href="style.css" rel="stylesheet">`+
			`<p>Hello, world!</p>`),
	)
	handlers.Handle(
		"example.org/style.css",
		stubTextHandler(`body { font-family: sans-serif; }`, "text/css"),
	)
	server := httptest.NewTLSServer(handlers)
	defer server.Close()

	dummyUA := "webpackager_test/0.1"
	header := make(http.Header)
	header.Set("User-Agent", dummyUA)

	config := makeConfig(server)
	config.RequestHeader = header
	pkg := webpackager.NewPackager(config)
	pkg.Run(urlutil.MustParse("https://example.org/hello.html"),
		exchange.NewValidPeriod(date, expires))

	for _, req := range pkg.FetchClient.(*fetchtest.FetchClient).Requests() {
		if got := req.Header.Get("User-Agent"); got != dummyUA {
			t.Errorf(`Requests()[%q].Header.Get("User-Agent") = %q, want %q`, req.URL, got, dummyUA)
		}
	}
}

func TestNoExchanges(t *testing.T) {
	handlers := http.NewServeMux()
	handlers.Handle(
		"example.org/hello.html",
		stubHTMLHandler(`<!doctype html><p>Hello, world!</p>`),
	)
	handlers.Handle(
		"example.org/secret.html",
		stubErrorHandler(http.StatusForbidden),
	)
	handlers.Handle(
		"example.org/redirect.html",
		http.RedirectHandler("hello.html", http.StatusFound),
	)
	server := httptest.NewTLSServer(handlers)
	defer server.Close()
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Redirected",
			url:  "https://example.org/redirect.html",
		}, {
			name: "NonOKStatus",
			url:  "https://example.org/secret.html",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := urlutil.MustParse(test.url)

			pkg := webpackager.NewPackager(makeConfig(server))
			pkg.Run(url, exchange.NewValidPeriod(date, expires))

			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			got, err := pkg.ResourceCache.Lookup(req)
			if err != nil {
				t.Fatalf("Lookup(%q) = error(%q), want success", url, err)
			}
			if got != nil {
				t.Errorf("Lookup(%q) = %v, want nil", url, got)
			}
		})
	}
}
