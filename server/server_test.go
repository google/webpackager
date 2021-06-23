// Copyright 2020 Google LLC
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

package server_test

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/fetch/fetchtest"
	"github.com/google/webpackager/internal/certchaintest"
	"github.com/google/webpackager/internal/timeutil"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/server"
	"github.com/google/webpackager/server/tomlconfig"
	"github.com/google/webpackager/urlmatcher"
	"github.com/google/webpackager/validity"
)

const cborFile = "../testdata/certs/cbor/ecdsap256_nosct.cbor"

func setupServer(www *httptest.Server) (*server.Server, string) {
	ac := certchaintest.MustReadAugmentedChainFile(cborFile)

	certManager := certmanager.NewManager(certmanager.Config{
		RawChainSource: &stubRawChainSource{ac.RawChain},
		OCSPRespSource: &stubOCSPRespSource{ac.OCSPResp},
		Cache:          newStubCache(),
	})

	s := server.NewServer(new(http.Server), server.Config{
		ServerConfig: tomlconfig.ServerConfig{
			DocPath:      "/priv/doc",
			CertPath:     "/webpkg/cert",
			ValidityPath: "/webpkg/validity",
			HealthPath:   "/healthz",
			SignParam:    "sign",
		},
		AllowTestCert: true,
		CertManager:   certManager,
		Packager: webpackager.NewPackager(webpackager.Config{
			FetchClient: fetch.WithSelector(
				fetchtest.NewFetchClient(www),
				&fetch.Selector{
					Allow: []urlmatcher.Matcher{
						urlmatcher.AllOf(
							urlmatcher.HasHostname("example.com"),
							urlmatcher.HasEscapedPathRegexp(regexp.MustCompile(`\A/public/.*\z`)),
						),
					},
				},
			),
			ValidityURLRule: validity.FixedURL(urlutil.MustParse("/webpkg/validity")),
			ExchangeFactory: server.NewExchangeMetaFactory(server.ExchangeConfig{
				CertManager: certManager,
				CertURLBase: urlutil.MustParse("/webpkg/cert"),
				PrivateKey:  certchaintest.MustReadPrivateKeyFile("../testdata/keys/ecdsap256.key"),
			}),
		}),
	})

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	go (func() { s.Serve(l) })()

	return s, l.Addr().String()
}

func setupContentServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/public/hello.html", func(w http.ResponseWriter, r *http.Request) {
		html := "<!doctype html><p>Hello, world!</p>"
		http.ServeContent(w, r, "hello.html", time.Time{}, strings.NewReader(html))
	})
	mux.HandleFunc("/public/page.cgi", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("id") {
		case "hello":
			html := "<!doctype html><p>Hello, world!</p>"
			http.ServeContent(w, r, "hello.html", time.Time{}, strings.NewReader(html))
		default:
			http.Error(w, "404 Not Found", http.StatusNotFound)
		}
	})
	mux.HandleFunc("/private/hello.html", func(w http.ResponseWriter, r *http.Request) {
		html := "<!doctype html><p>hello, world</p>"
		http.ServeContent(w, r, "hello.html", time.Time{}, strings.NewReader(html))
	})

	return httptest.NewTLSServer(mux)
}

func TestHandleDoc(t *testing.T) {
	www := setupContentServer()
	defer www.Close()
	s, addr := setupServer(www)
	defer s.Close()

	tests := []struct {
		name   string
		url    string
		accept string
	}{
		{
			name:   "EmbeddedURL",
			url:    "http://" + addr + "/priv/doc/https://example.com/public/hello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "EmbeddedURLWithQuery",
			url:    "http://" + addr + "/priv/doc/https://example.com/public/page.cgi?id=hello",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignParam",
			url:    "http://" + addr + "/priv/doc?sign=https%3A%2F%2Fexample.com%2Fpublic%2Fhello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "ComplexAcceptHeader",
			url:    "http://" + addr + "/priv/doc/https://example.com/public/hello.html",
			accept: "*/*;q=0.1,application/signed-exchange;q=1.0;v=b3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			timeutil.StubNowToAdjust(time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC))

			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Accept", test.accept)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			wantType := "application/signed-exchange;v=b3"

			if got := resp.StatusCode; got != http.StatusOK {
				t.Errorf("StatusCode = %v, want %v", got, http.StatusOK)
			}
			if got := resp.Header.Get("Content-Type"); got != wantType {
				t.Errorf("[Content-Type] = %q, want %q", got, wantType)
			}
			if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("[X-Content-Type-Options] = %q, want %q", got, "nosniff")
			}

			sxg, err := signedexchange.ReadExchange(resp.Body)
			if err != nil {
				t.Errorf("ReadExchange() = error(%q), want success", err)
			}
			_, ok := sxg.Verify(
				timeutil.Now(),
				func(url string) ([]byte, error) { return ioutil.ReadFile(cborFile) },
				log.New(os.Stderr, "", log.LstdFlags),
			)
			if !ok {
				t.Fatalf("Verify() = !ok, want ok")
			}
		})
	}
}

func TestHandleDoc_ClientError(t *testing.T) {
	www := setupContentServer()
	defer www.Close()
	s, addr := setupServer(www)
	defer s.Close()

	tests := []struct {
		name   string
		url    string
		accept string
	}{
		{
			name:   "Accept_MissingSXG",
			url:    "http://" + addr + "/priv/doc/https://example.com/public/hello.html",
			accept: "text/html,application/xhtml+xml;q=0.9",
		},
		{
			name:   "SignURL_Missing",
			url:    "http://" + addr + "/priv/doc/",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignURL_UnknownParam",
			url:    "http://" + addr + "/priv/doc?ngis=https%3A%2F%2Fexample.com%2Fpublic%2Fhello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignURL_NotAbsolute",
			url:    "http://" + addr + "/priv/doc?sign=%2Fpublic%2Fhello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignURL_NotParsable",
			url:    "http://" + addr + "/priv/doc/http$://example.com/public/hello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignURL_NotHTTPS",
			url:    "http://" + addr + "/priv/doc?sign=http%3A%2F%2Fexample.com%2Fpublic%2Fhello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignURL_HasUserInfo",
			url:    "http://" + addr + "/priv/doc?sign=https%3A%2F%2Fuser%3Apass%40example.com%2Fpublic%2Fhello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "SignURL_HasFragment",
			url:    "http://" + addr + "/priv/doc?sign=https%3A%2F%2Fexample.com%2Fpublic%2Fhello.html%23hash",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "Domain_NotForFetch",
			url:    "http://" + addr + "/priv/doc/https://example.org/public/hello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "Domain_TrickAttempted",
			url:    "http://" + addr + "/priv/doc/https://example.org/../example.com/public/hello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "Path_NotForFetch",
			url:    "http://" + addr + "/priv/doc/https://example.com/private/hello.html",
			accept: "application/signed-exchange;v=b3",
		},
		{
			name:   "Path_TrickAttempted",
			url:    "http://" + addr + "/priv/doc/https://example.com/public/%2E%2E/private/hello.html",
			accept: "application/signed-exchange;v=b3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Accept", test.accept)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if got := resp.StatusCode; got != http.StatusBadRequest {
				t.Errorf("StatusCode = %v, want %v", got, http.StatusBadRequest)
			}
		})
	}
}

func TestHandleCert(t *testing.T) {
	www := setupContentServer()
	defer www.Close()
	s, addr := setupServer(www)
	defer s.Close()

	// TODO(banaag): Insert 5 second delay to test flake theory.
	time.Sleep(5 * time.Second)

	url := "http://" + addr + "/webpkg/cert/qwk4hz4Swff9wKMvr1hri3YH4MeFAH8_PE9jnJ9nx6A"

	wantBody, err := ioutil.ReadFile(cborFile)
	if err != nil {
		t.Fatal(err)
	}
	wantType := "application/cert-chain+cbor"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", got, http.StatusOK)
	}
	if got := resp.Header.Get("Content-Type"); got != wantType {
		t.Errorf("[Content-Type] = %q, want %q", got, wantType)
	}
	if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("[X-Content-Type-Options] = %q, want %q", got, "nosniff")
	}

	gotBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantBody, gotBody); diff != "" {
		t.Errorf("Body mismatch (-want +got):\n%s", diff)
	}
}

func TestHandlerCert_DigestMismatch(t *testing.T) {
	www := setupContentServer()
	defer www.Close()
	s, addr := setupServer(www)
	defer s.Close()

	url := "http://" + addr + "/webpkg/cert/k8HZqkHWuFLy34Bc0R-QKD0Vkb7LwoM_ckBc_li0Nzc"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got := resp.StatusCode; got != http.StatusNotFound {
		t.Errorf("StatusCode = %v, want %v", got, http.StatusNotFound)
	}
}

func TestHandleHealth(t *testing.T) {
	www := setupContentServer()
	defer www.Close()
	s, addr := setupServer(www)
	defer s.Close()

	url := "http://" + addr + "/healthz"
	wantBody := []uint8{0x6f, 0x6b} // "ok"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", got, http.StatusOK)
	}

	gotBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantBody, gotBody); diff != "" {
		t.Errorf("Body mismatch (-want +got):\n%s", diff)
	}
}

func TestHandleValidity(t *testing.T) {
	www := setupContentServer()
	defer www.Close()
	s, addr := setupServer(www)
	defer s.Close()

	url := "http://" + addr + "/webpkg/validity"

	wantBody := []byte{0xa0}
	wantType := "application/cbor"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", got, http.StatusOK)
	}
	if got := resp.Header.Get("Content-Type"); got != wantType {
		t.Errorf("[Content-Type] = %q, want %q", got, wantType)
	}
	if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("[X-Content-Type-Options] = %q, want %q", got, "nosniff")
	}

	gotBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantBody, gotBody); diff != "" {
		t.Errorf("Body mismatch (-want +got):\n%s", diff)
	}
}
