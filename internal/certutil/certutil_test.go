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

package certutil

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFakeCertFetcher(t *testing.T) {
	certChain, err := ReadCertChainFile("../../testdata/certs/test.cbor")
	if err != nil {
		t.Fatal(err)
	}
	certURL, err := url.Parse("https://example.org/cert.cbor")
	if err != nil {
		t.Fatal(err)
	}
	fakeCertFetcher := FakeCertFetcher(certChain, certURL)

	t.Run("Success", func(t *testing.T) {
		want, err := ioutil.ReadFile("../../testdata/certs/test.cbor")
		if err != nil {
			t.Fatal(err)
		}

		got, err := fakeCertFetcher("https://example.org/cert.cbor")
		if err != nil {
			t.Fatalf("got error(%q), want success", err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("got %q (%d bytes), want %q (%d bytes)", got, len(got), want, len(want))
		}
	})

	t.Run("UnknownURL", func(t *testing.T) {
		_, err := fakeCertFetcher("https://example.org/test.cbor")
		if err != errUnknownCertURL {
			if err != nil {
				t.Errorf("got error(%q), want error(%q)", err, errUnknownCertURL)
			} else {
				t.Errorf("got success, want error(%q)", errUnknownCertURL)
			}
		}
	})
}

func TestGetCertificates(t *testing.T) {
	tests := []struct {
		name string
		cbor string
		pem  string
	}{
		{
			name: "test.cbor",
			cbor: "../../testdata/certs/test.cbor",
			pem:  "../../testdata/certs/test.pem",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cbor, err := ReadCertChainFile(test.cbor)
			if err != nil {
				t.Fatal(err)
			}
			want, err := ReadCertificatesFile(test.pem)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, GetCertificates(cbor)); diff != "" {
				t.Errorf("GetCertificates() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReadPrivateKeyFile_Success(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "test.key",
			file: "../../testdata/certs/test.key",
		},
		{
			name: "pkcs8-ecdsa.key",
			file: "../../testdata/certs/pkcs8-ecdsa.key",
		},
		{
			name: "pkcs8-multi.key",
			file: "../../testdata/certs/pkcs8-multi.key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ReadPrivateKeyFile(test.file); err != nil {
				t.Errorf("got error(%q), want success", err)
			}
		})
	}
}

func TestReadPrivateKeyFile_Error(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "pkcs8-rsa.key",
			file: "../../testdata/certs/pkcs8-rsa.key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ReadPrivateKeyFile(test.file); err == nil {
				t.Error("got success, want error")
			}
		})
	}
}
