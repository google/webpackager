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

package exchange_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/WICG/webpackage/go/signedexchange/structuredheader"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/internal/certchaintest"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/preload"
)

func eraseSignature(sxg []byte) []byte {
	re := regexp.MustCompile(`;sig=\*([A-Za-z0-9+/]*=*)\*`)
	return re.ReplaceAll(sxg, []byte(";sig=*/erased/*"))
}

func TestFactory(t *testing.T) {
	// Initialize Factory with the parameters given to gen-signedexchange.
	factory := exchange.NewFactory(exchange.Config{
		CertChain:  certchaintest.MustReadAugmentedChainFile("../testdata/certs/cbor/ecdsap256_nosct.cbor"),
		CertURL:    urlutil.MustParse("https://example.org/cert.cbor"),
		PrivateKey: certchaintest.MustReadPrivateKeyFile("../testdata/keys/ecdsap256.key"),
	})
	vp := exchange.NewValidPeriod(
		time.Date(2019, time.April, 22, 19, 30, 0, 0, time.UTC),
		time.Date(2019, time.April, 29, 19, 30, 0, 0, time.UTC))
	// Create Response with the content passed to gen-signedexchange.
	tests := []struct {
		name     string
		url      string
		respHead string
		htmlFile string
		preloads []*preload.Preload
		validity string
		sxgFile  string
	}{
		{
			name: "standalone.sxg",
			url:  "https://example.org/standalone.html",
			respHead: "HTTP/1.1 200 OK\r\n" +
				"Cache-Control: public, max-age=604800\r\n" +
				"Content-Length: 37\r\n" +
				"Content-Type: text/html;charset=utf-8\r\n",
			htmlFile: "../testdata/sxg/standalone.html",
			preloads: nil,
			validity: "https://example.org/standalone.html.validity.1555961400",
			sxgFile:  "../testdata/sxg/standalone.sxg",
		},
		{
			name: "preloading.sxg",
			url:  "https://example.org/preloading.html",
			respHead: "HTTP/1.1 200 OK\r\n" +
				"Cache-Control: public, max-age=604800\r\n" +
				"Content-Length: 78\r\n" +
				"Content-Type: text/html;charset=utf-8\r\n",
			htmlFile: "../testdata/sxg/preloading.html",
			preloads: []*preload.Preload{
				preload.NewPreloadForResource(
					&resource.Resource{
						RequestURL: urlutil.MustParse("https://example.org/style.css"),
						Integrity:  "dummy-integrity",
					},
					preload.AsStyle,
				),
			},
			validity: "https://example.org/preloading.html.validity.1555961400",
			sxgFile:  "../testdata/sxg/preloading.sxg",
		},
		{
			name: "incomplete.sxg",
			url:  "https://example.org/incomplete.html",
			respHead: "HTTP/1.1 200 OK\r\n" +
				"Cache-Control: public, max-age=604800\r\n" +
				"Content-Length: 78\r\n" +
				"Content-Type: text/html;charset=utf-8\r\n",
			htmlFile: "../testdata/sxg/incomplete.html",
			preloads: []*preload.Preload{
				preload.NewPreloadForResource(
					&resource.Resource{
						RequestURL: urlutil.MustParse("https://example.org/style.css"),
						Integrity:  "", // Missing
					},
					preload.AsStyle,
				),
			},
			validity: "https://example.org/incomplete.html.validity.1555961400",
			sxgFile:  "../testdata/sxg/incomplete.sxg",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			html, err := ioutil.ReadFile(test.htmlFile)
			if err != nil {
				t.Fatal(err)
			}
			resp := exchangetest.MakeResponse(
				test.url, fmt.Sprintf("%s\r\n%s", test.respHead, html))
			resp.Preloads = test.preloads
			e, err := factory.NewExchange(resp, vp, urlutil.MustParse(test.validity))
			if err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			// Check the byte sequence matches the pre-generated file.
			var b bytes.Buffer
			if err := e.Write(&b); err != nil {
				t.Errorf("got error(%q), want success", err)
			} else {
				got := b.Bytes()
				want, err := ioutil.ReadFile(test.sxgFile)
				if err != nil {
					t.Fatal(err)
				}
				// Signatures change every time they are generated, so exclude.
				// The signature is tested below with Verify().
				got = eraseSignature(got)
				want = eraseSignature(want)
				if !bytes.Equal(got, want) {
					t.Errorf("got %q (%d bytes), want %q (%d bytes)", got, len(got), want, len(want))
				}
			}
			// Check Verify succeeds and returns the correct payload.
			if got, err := factory.Verify(e, vp.Date()); err != nil {
				t.Errorf("got error(%q), want success", err)
			} else {
				if !bytes.Equal(got, html) {
					t.Errorf("got %q (%d bytes), want %q (%d bytes)", got, len(got), html, len(html))
				}
			}
		})
	}
}

func TestRelativeCertURL(t *testing.T) {
	factory := exchange.NewFactory(exchange.Config{
		CertChain:  certchaintest.MustReadAugmentedChainFile("../testdata/certs/cbor/ecdsap256_nosct.cbor"),
		CertURL:    urlutil.MustParse("/cert.cbor"),
		PrivateKey: certchaintest.MustReadPrivateKeyFile("../testdata/keys/ecdsap256.key"),
	})
	vp := exchange.NewValidPeriod(
		time.Date(2019, time.April, 22, 19, 30, 0, 0, time.UTC),
		time.Date(2019, time.April, 29, 19, 30, 0, 0, time.UTC))

	tests := []struct {
		name        string
		url         string
		wantCertURL string
	}{
		{
			name:        "example.org",
			url:         "https://example.org/index.html",
			wantCertURL: "https://example.org/cert.cbor",
		},
		{
			name:        "example.com",
			url:         "https://example.com/index.html",
			wantCertURL: "https://example.com/cert.cbor",
		},
	}

	vu := urlutil.MustParse("https://example.org/index.html.validity")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeEmptyResponse(test.url)
			e, err := factory.NewExchange(resp, vp, vu)
			if err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			sig, err := structuredheader.ParseParameterisedList(e.SignatureHeaderValue)
			if err != nil {
				t.Fatalf("ParseParameterizedList() = error(%q), want success", err)
			}
			if len(sig) == 0 {
				t.Fatal("ParseParameterizedList() = empty, want nonempty")
			}
			if got := sig[0].Params["cert-url"]; got != test.wantCertURL {
				t.Errorf(`sig[0].Params["cert-url"] = %q, want %q`, got, test.wantCertURL)
			}
		})
	}
}
