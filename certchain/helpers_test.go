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
	"crypto/x509"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certchainutil"
)

// bigIntComparer allows big.Int values to be compared in the cmp methods.
// big.Int provides Cmp(), but not Equal().
var bigIntComparer = cmp.Comparer(func(x, y *big.Int) bool {
	return x.Cmp(y) == 0
})

func mustReadCertificates(filename string) []*x509.Certificate {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	certs, err := signedexchange.ParseCertificates(pemBytes)
	if err != nil {
		panic(err)
	}
	return certs
}

func mustReadRawChainFile(t *testing.T, filename string) *certchain.RawChain {
	t.Helper()

	got, err := certchainutil.ReadRawChainFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func mustReadOCSPRespFile(t *testing.T, filename string) *certchain.OCSPResponse {
	t.Helper()

	got, err := certchainutil.ReadOCSPRespFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func mustReadAugmentedChainFile(t *testing.T, filename string) *certchain.AugmentedChain {
	t.Helper()

	got, err := certchainutil.ReadAugmentedChainFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return got
}
