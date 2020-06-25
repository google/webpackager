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
	"crypto/x509"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
)

var (
	rootCA     = mustReadCertificates("../testdata/CA/root/cert.pem")[0]
	interCA    = mustReadCertificates("../testdata/CA/inter/cert.pem")[0]
	issuedCert = mustReadCertificates("../testdata/certs/issued/ecdsap256_sxg_60days.crt")[0]
	selfCert   = mustReadCertificates("../testdata/certs/chain/self_signed.pem")[0]
)

func TestNewRawChain_Success(t *testing.T) {
	tests := []struct {
		name  string
		certs []*x509.Certificate
		want  *certchain.RawChain
	}{
		{
			name:  "BaseCase",
			certs: []*x509.Certificate{issuedCert, interCA, rootCA},
			want: &certchain.RawChain{
				Certs:      []*x509.Certificate{issuedCert, interCA, rootCA},
				Digest:     "qwk4hz4Swff9wKMvr1hri3YH4MeFAH8_PE9jnJ9nx6A",
				Leaf:       issuedCert,
				Issuer:     interCA,
				OCSPServer: "http://ocsp.webpackager-ca.test",
			},
		},
		{
			name:  "WithoutRoot",
			certs: []*x509.Certificate{issuedCert, interCA},
			want: &certchain.RawChain{
				Certs:      []*x509.Certificate{issuedCert, interCA},
				Digest:     "ZonVSpXPO9OROUKwvjw8kDzYp_RmJyiImr3g-dlMtAw",
				Leaf:       issuedCert,
				Issuer:     interCA,
				OCSPServer: "http://ocsp.webpackager-ca.test",
			},
		},
		{
			name:  "SelfSigned",
			certs: []*x509.Certificate{selfCert},
			want: &certchain.RawChain{
				Certs:      []*x509.Certificate{selfCert},
				Digest:     "k8HZqkHWuFLy34Bc0R-QKD0Vkb7LwoM_ckBc_li0Nzc",
				Leaf:       selfCert,
				Issuer:     selfCert,
				OCSPServer: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := certchain.NewRawChain(test.certs)
			if err != nil {
				t.Fatalf("NewRawChain() = error(%q), want success", err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("RawChain mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewRawChain_Error(t *testing.T) {
	tests := []struct {
		name  string
		certs []*x509.Certificate
	}{
		{
			name:  "Empty",
			certs: nil,
		},
		{
			name:  "MissingIssuer",
			certs: []*x509.Certificate{issuedCert, rootCA},
		},
		{
			name:  "WrongOrder",
			certs: []*x509.Certificate{issuedCert, rootCA, interCA},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := certchain.NewRawChain(test.certs)
			if err == nil {
				t.Fatal("NewRawChain() = success, want error")
			}
		})
	}
}

func TestWritePEM(t *testing.T) {
	c, err := certchain.NewRawChain([]*x509.Certificate{
		issuedCert,
		interCA,
		rootCA,
	})
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("../testdata/certs/chain/ecdsap256.pem")
	if err != nil {
		t.Fatal(err)
	}

	var got bytes.Buffer

	if err := c.WritePEM(&got); err != nil {
		t.Errorf("WritePEM() = error(%q), want success", err)
	}
	if diff := cmp.Diff(want, got.Bytes()); diff != "" {
		t.Errorf("PEM mismatch (-want +got):\n%s", diff)
	}
}

func TestVerifyChain_Success(t *testing.T) {
	tests := []struct {
		name     string
		rawChain *certchain.RawChain
		time     time.Time
	}{
		{
			name:     "OnMidpoint",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "OnNotBefore",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "OnNotAfter",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "MissingRoot",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/without_root.pem"),
			time:     time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "SelfSigned",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/self_signed.pem"),
			time:     time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.rawChain.VerifyChain(test.time); err != nil {
				t.Errorf("VerifyChain() = error(%q), want success", err)
			}
		})
	}
}

func TestVerifyChain_Error(t *testing.T) {
	tests := []struct {
		name     string
		rawChain *certchain.RawChain
		time     time.Time
	}{
		{
			name:     "BeforePeriod",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.March, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "AfterPeriod",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
			time:     time.Date(2020, time.May, 31, 0, 0, 1, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.rawChain.VerifyChain(test.time); err == nil {
				t.Errorf("VerifyChain() = nil, want error")
			}
		})
	}
}

func TestVerifySXGCriteria_Success(t *testing.T) {
	tests := []struct {
		name     string
		rawChain *certchain.RawChain
	}{
		{
			name:     "BaseCase",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap256.pem"),
		},
		{
			name:     "AlternativePublicKey",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap384.pem"),
		},
		{
			name:     "MaxCertDuration",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/lasting_90days.pem"),
		},
		{
			name:     "SelfSigned",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/self_signed.pem"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.rawChain.VerifySXGCriteria(); err != nil {
				t.Errorf("VerifySXGCriteria() = error(%q), want success", err)
			}
		})
	}
}

func TestVerifySXGCriteria_Error(t *testing.T) {
	tests := []struct {
		name     string
		rawChain *certchain.RawChain
	}{
		{
			name:     "BadPublicKey_P521IsUnsupported",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/ecdsap521.pem"),
		},
		{
			name:     "BadPublicKey_RSAIsInvalid",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/rsa4096.pem"),
		},
		{
			name:     "BadCertDuration_Negative",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/lasting_-1days.pem"),
		},
		{
			name:     "BadCertDuration_OffByOne",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/lasting_91days.pem"),
		},
		{
			name:     "BadCertDuration_ClearlyTooLong",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/lasting_365days.pem"),
		},
		{
			name:     "MissingCanSignHttpExchanges",
			rawChain: mustReadRawChainFile(t, "../testdata/certs/chain/non_sxg_cert.pem"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.rawChain.VerifySXGCriteria(); err == nil {
				t.Errorf("VerifySXGCriteria() = success, want error")
			}
		})
	}
}
