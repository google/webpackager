// Copyright 2020 Google LLC
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

package certchain_test

import (
	"bytes"
	"crypto/x509/pkix"
	"encoding/asn1"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
)

func TestReadAugmentedChain_Issued(t *testing.T) {
	cbor, err := ioutil.ReadFile("../testdata/certs/cbor/ecdsap256_nosct.cbor")
	if err != nil {
		t.Fatal(err)
	}

	got, err := certchain.ReadAugmentedChain(bytes.NewReader(cbor))

	if err != nil {
		t.Fatalf("ReadAugmentedChain() = error(%q), want success", err)
	}

	wantRawChain := mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem")
	if diff := cmp.Diff(wantRawChain, got.RawChain); diff != "" {
		t.Errorf("RawChain mismatch (-want +got):\n%s", diff)
	}
	wantOCSPResp := mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp")
	if diff := cmp.Diff(wantOCSPResp, got.OCSPResp, bigIntComparer); diff != "" {
		t.Errorf("OCSPResp mismatch (-want +got):\n%s", diff)
	}
}

func TestReadAugmentedChain_SelfSigned(t *testing.T) {
	cbor, err := ioutil.ReadFile("../testdata/certs/cbor/self_signed.cbor")
	if err != nil {
		t.Fatal(err)
	}

	got, err := certchain.ReadAugmentedChain(bytes.NewReader(cbor))

	if err != certchain.ErrInvalidOCSPValue {
		if err != nil {
			t.Fatalf("ReadAugmentedChain() = error(%q), want ErrInvalidOCSPValue", err)
		} else {
			t.Errorf("ReadAugmentedChain() = success, want ErrInvalidOCSPValue")
		}
	}

	wantRawChain := mustReadRawChainFile(t, "../testdata/certs/chain/self_signed.pem")
	if diff := cmp.Diff(wantRawChain, got.RawChain); diff != "" {
		t.Errorf("RawChain mismatch (-want +got):\n%s", diff)
	}
	wantOCSPResp := certchain.DummyOCSPResponse
	if diff := cmp.Diff(wantOCSPResp, got.OCSPResp, bigIntComparer); diff != "" {
		t.Errorf("OCSPResp mismatch (-want +got):\n%s", diff)
	}
}

func TestWriteCBOR(t *testing.T) {
	ac := certchain.NewAugmentedChain(
		mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
		mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
		nil, // sct
	)

	want, err := ioutil.ReadFile("../testdata/certs/cbor/ecdsap256_nosct.cbor")
	if err != nil {
		t.Fatal(err)
	}
	var got bytes.Buffer
	if err := ac.WriteCBOR(&got); err != nil {
		t.Fatalf("WriteCBOR() = error(%q), want success", err)
	}
	if diff := cmp.Diff(want, got.Bytes()); diff != "" {
		t.Errorf("CBOR mismatch (-want +got):\n%s", diff)
	}
}

func TestHasSCTList(t *testing.T) {
	// TODO(yuizumi): Use real certificates; stop mutating the parsed ones.

	tests := []struct {
		name string
		ac   *certchain.AugmentedChain
		want bool
	}{
		{
			name: "EmbeddedInCert",
			ac: (func() *certchain.AugmentedChain {
				ac := mustReadAugmentedChainFile(t, "../testdata/certs/cbor/ecdsap256_nosct.cbor")
				ac.Certs[0].Extensions = []pkix.Extension{
					pkix.Extension{Id: asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 11129, 2, 4, 2}},
				}
				return ac
			})(),
			want: true,
		},
		{
			name: "EmbeddedInOCSP",
			ac: (func() *certchain.AugmentedChain {
				ac := mustReadAugmentedChainFile(t, "../testdata/certs/cbor/ecdsap256_nosct.cbor")
				ac.OCSPResp.Extensions = []pkix.Extension{
					pkix.Extension{Id: asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 11129, 2, 4, 5}},
				}
				return ac
			})(),
			want: true,
		},
		{
			name: "Unembedded",
			ac: (func() *certchain.AugmentedChain {
				ac := mustReadAugmentedChainFile(t, "../testdata/certs/cbor/ecdsap256_nosct.cbor")
				ac.SCTList = []byte("dummy-sct")
				return ac
			})(),
			want: true,
		},
		{
			name: "Missing",
			ac: (func() *certchain.AugmentedChain {
				ac := mustReadAugmentedChainFile(t, "../testdata/certs/cbor/ecdsap256_nosct.cbor")
				return ac
			})(),
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.ac.HasSCTList(); got != test.want {
				t.Errorf("HasSCTList() = %v, want %v", got, test.want)
			}
		})
	}
}
