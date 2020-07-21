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

package acmeclient_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-acme/lego/v3/acme"
	"github.com/go-acme/lego/v3/platform/tester"
	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/certchain/certmanager/acmeclient"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"github.com/google/webpackager/internal/certchaintest"
	jose "gopkg.in/square/go-jose.v2"
)

var fakeACMECert = certchaintest.MustReadRawChainFile("../../../testdata/certs/chain/fake_acme_cert.pem")
var testCert = certchaintest.MustReadRawChainFile("../../../testdata/certs/chain/ecdsap256.pem")

var fakeCSR = x509.CertificateRequest{
	Subject: pkix.Name{
		CommonName:   "test.example.com",
		Organization: []string{"Acme Co"},
	},
	DNSNames: []string{"test.example.com"},
}

func TestFetchSuccess(t *testing.T) {
	mux, apiURL, tearDown := tester.SetupFakeAPI()
	defer tearDown()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() = error(%q), want success", err)
	}

	setupMux(mux, apiURL, privateKey)
	mux.HandleFunc("/certificate", func(w http.ResponseWriter, _ *http.Request) {
		fakeACMECert.WritePEM(w)
	})

	config := acmeclient.Config{
		CertSignRequest:   &fakeCSR,
		User:              acmeclient.NewUser("test@example.com", privateKey),
		DiscoveryURL:      apiURL + "/dir",
		HTTPChallengePort: getFreeTCPPort(t),
		FetchTiming:       certmanager.FetchOnlyOnce(),
	}
	client, err := acmeclient.NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() = error(%q), want success", err)
	}

	now := time.Date(2020, time.May, 1, 15, 0, 0, 0, time.UTC)
	chain, nextRun, err := client.Fetch(nil, func() time.Time { return now })

	// nextRun must be valid even if (err != nil).
	if _, ok := nextRun.(*futureevent.NeverOccursEvent); !ok {
		t.Errorf("nextRun = %#v, want NeverOccursEvent", nextRun)
	}
	if err != nil {
		t.Fatalf("client.Fetch() = error(%q), want success", err)
	}
	if diff := cmp.Diff(fakeACMECert, chain); diff != "" {
		t.Errorf("client.Fetch() mismatch (-want +got):\n%s", diff)
	}

	now = time.Date(2020, time.May, 1, 15, 0, 0, 0, time.UTC)
	chain, nextRun, err = client.Fetch(testCert, func() time.Time { return now })

	// nextRun must be valid even if (err != nil).
	if _, ok := nextRun.(*futureevent.NeverOccursEvent); !ok {
		t.Errorf("nextRun = %#v, want NeverOccursEvent", nextRun)
	}
	if err != nil {
		t.Fatalf("client.Fetch() = error(%q), want success", err)
	}
	if diff := cmp.Diff(testCert, chain); diff != "" {
		t.Errorf("client.Fetch() mismatch (-want +got):\n%s", diff)
	}
}

func TestFetchFailure(t *testing.T) {
	mux, apiURL, tearDown := tester.SetupFakeAPI()
	defer tearDown()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() = error(%q), want success", err)
	}

	setupMux(mux, apiURL, privateKey)
	mux.HandleFunc("/certificate", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "", http.StatusInternalServerError)
	})

	config := acmeclient.Config{
		CertSignRequest:   &fakeCSR,
		User:              acmeclient.NewUser("test@example.com", privateKey),
		DiscoveryURL:      apiURL + "/dir",
		HTTPChallengePort: getFreeTCPPort(t),
		FetchTiming:       certmanager.FetchOnlyOnce(),
	}
	client, err := acmeclient.NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() = error(%q), want success", err)
	}

	now := time.Date(2020, time.May, 1, 15, 0, 0, 0, time.UTC)
	chain, nextRun, err := client.Fetch(nil, func() time.Time { return now })

	// nextRun must be valid even if (err != nil).
	if _, ok := nextRun.(*futureevent.NeverOccursEvent); !ok {
		t.Errorf("nextRun = %#v, want NeverOccursEvent", nextRun)
	}
	if err == nil {
		t.Fatalf("client.Fetch() = %#v, want error", chain)
	}
}

func setupMux(mux *http.ServeMux, apiURL string, privateKey *rsa.PrivateKey) {
	mux.HandleFunc("/newOrder", func(w http.ResponseWriter, req *http.Request) {
		handleBody(w, req, apiURL, privateKey)
	})
	mux.HandleFunc("/finalize", func(w http.ResponseWriter, req *http.Request) {
		handleBody(w, req, apiURL, privateKey)
	})
}

func handleBody(w http.ResponseWriter, req *http.Request, apiURL string, privateKey *rsa.PrivateKey) {
	if req.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jws, err := jose.ParseSigned(string(reqBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := jws.Verify(&jose.JSONWebKey{
		Key:       privateKey.Public(),
		Algorithm: "RSA",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order acme.Order
	if err := json.Unmarshal(body, &order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.URL.Path {
	case "/newOrder":
		err = tester.WriteJSONResponse(w, acme.Order{
			Status:      acme.StatusValid,
			Identifiers: order.Identifiers,
			Finalize:    apiURL + "/finalize",
		})
	case "/finalize":
		err = tester.WriteJSONResponse(w, acme.Order{
			Status:      acme.StatusValid,
			Identifiers: order.Identifiers,
			Certificate: apiURL + "/certificate",
		})
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
