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
	"crypto"
	"crypto/x509"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/webpackager/certchain"
	"golang.org/x/crypto/ocsp"
)

func TestParseOCSPResponseForRawChain_Success(t *testing.T) {
	tests := []struct {
		name     string
		ocspFile string
		rawChain *certchain.RawChain
		want     *ocsp.Response
	}{
		{
			name:     "BaseCase",
			ocspFile: "../testdata/ocsp/ecdsap256_7days.ocsp",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			want: &ocsp.Response{
				Status:             ocsp.Good,
				SerialNumber:       big.NewInt(1),
				ProducedAt:         time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
				ThisUpdate:         time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
				NextUpdate:         time.Date(2020, time.May, 8, 0, 0, 0, 0, time.UTC),
				SignatureAlgorithm: x509.ECDSAWithSHA256,
				IssuerHash:         crypto.SHA1,
			},
		},
		{
			name:     "AnotherCert",
			ocspFile: "../testdata/ocsp/ecdsap384_7days.ocsp",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap384.pem"),
			want: &ocsp.Response{
				Status:             ocsp.Good,
				SerialNumber:       big.NewInt(2),
				ProducedAt:         time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
				ThisUpdate:         time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
				NextUpdate:         time.Date(2020, time.May, 8, 0, 0, 0, 0, time.UTC),
				SignatureAlgorithm: x509.ECDSAWithSHA256,
				IssuerHash:         crypto.SHA1,
			},
		},
	}

	cmpOCSPResponse := cmp.Options{
		bigIntComparer,                // Compare big.Int using Cmp().
		cmpopts.IgnoreTypes([]byte{}), // Ignore all []byte fields.
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			derBytes, err := ioutil.ReadFile(test.ocspFile)
			if err != nil {
				t.Fatal(err)
			}
			got, err := certchain.ParseOCSPResponseForRawChain(derBytes, test.rawChain)
			if err != nil {
				t.Errorf("ParseOCSPResponseForRawChain() = error(%q), want success", err)
			}
			if diff := cmp.Diff(test.want, got.Response, cmpOCSPResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(derBytes, got.Raw); diff != "" {
				t.Errorf("Raw mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseOCSPResponseForRawChain_Error(t *testing.T) {
	tests := []struct {
		name     string
		ocspFile string
		rawChain *certchain.RawChain
	}{
		{
			name:     "SerialMismatch",
			ocspFile: "../testdata/ocsp/ecdsap256_7days.ocsp",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap384.pem"),
		},
		{
			name:     "IssuerMismatch",
			ocspFile: "../testdata/ocsp/ecdsap256_7days.ocsp",
			rawChain: mustReadRawChainFile(t, "../testdata/CA/root/cert.pem"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			derBytes, err := ioutil.ReadFile(test.ocspFile)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := certchain.ParseOCSPResponseForRawChain(derBytes, test.rawChain); err == nil {
				t.Error("ParseOCSPResponseForRawChain() = success, want error")
			}
		})
	}
}

func TestOCSPResponseVerifyForRawChain_Success(t *testing.T) {
	tests := []struct {
		name     string
		ocspResp *certchain.OCSPResponse
		rawChain *certchain.RawChain
		time     time.Time
	}{
		{
			name:     "OnMidpoint",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 4, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "OnThisUpdate",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "OnNextUpdate",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "AnotherCert",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap384_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap384.pem"),
			time:     time.Date(2020, time.May, 8, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.ocspResp.VerifyForRawChain(test.time, test.rawChain); err != nil {
				t.Errorf("VerifyForRawChain() = error(%q), want success", err)
			}
		})
	}
}

func TestOCSPResponseVerifyForRawChain_Error(t *testing.T) {
	tests := []struct {
		name     string
		ocspResp *certchain.OCSPResponse
		rawChain *certchain.RawChain
		time     time.Time
	}{
		{
			name:     "BeforePeriod",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.April, 30, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "AfterPeriod",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 8, 0, 0, 1, 0, time.UTC),
		},
		{
			name:     "SerialMismatch",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap384.pem"),
			time:     time.Date(2020, time.May, 4, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "IssuerMismatch",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
			rawChain: mustReadRawChainFile(t, "../testdata/CA/root/cert.pem"),
			time:     time.Date(2020, time.May, 4, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.ocspResp.VerifyForRawChain(test.time, test.rawChain); err == nil {
				t.Errorf("VerifyForRawChain() = success, want error")
			}
		})
	}
}

func TestOCSPResponseVerifyForRawChain_Dummy(t *testing.T) {
	ocspResp := certchain.DummyOCSPResponse
	rawChain := mustReadRawChainFile(t, "../testdata/certs/chain/self_signed.pem")
	time := time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC)

	if err := ocspResp.VerifyForRawChain(time, rawChain); err != certchain.ErrDummyOCSPResponse {
		t.Errorf("VerifyForRawChain() = error(%q), want ErrDummyOCSPResponse", err)
	}
}

func TestOCSPResponseVerifyForOCSPDummy(t *testing.T) {
	ocspResp := certchain.DummyOCSPResponse

	if _, err := certchain.ParseOCSPResponse(ocspResp.Raw); err != nil {
		t.Errorf("ParseOCSPResponse() = error(%q), want success", err)
	}
}
func TestOCSPResponseVerifySXGCriteria_Success(t *testing.T) {
	tests := []struct {
		name     string
		ocspResp *certchain.OCSPResponse
	}{
		{
			name:     "BaseCase",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_7days.ocsp"),
		},
		{
			name:     "AnotherCert",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap384_7days.ocsp"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.ocspResp.VerifySXGCriteria(); err != nil {
				t.Errorf("VerifySXGCriteria() = error(%q), want success", err)
			}
		})
	}
}

func TestOCSPResponseVerifySXGCriteria_Error(t *testing.T) {
	tests := []struct {
		name     string
		ocspResp *certchain.OCSPResponse
	}{
		{
			name:     "RevokedCert",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/revoked_7days.ocsp"),
		},
		{
			name:     "BadDuration",
			ocspResp: mustReadOCSPRespFile(t, "../testdata/ocsp/ecdsap256_8days.ocsp"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.ocspResp.VerifySXGCriteria(); err == nil {
				t.Errorf("VerifySXGCriteria() = success, want error")
			}
		})
	}
}

func TestOCSPResponseVerifySXGCriteria_Dummy(t *testing.T) {
	ocspResp := certchain.DummyOCSPResponse

	if err := ocspResp.VerifySXGCriteria(); err != certchain.ErrDummyOCSPResponse {
		t.Errorf("VerifyForRawChain() = error(%q), want ErrDummyOCSPResponse", err)
	}
}
