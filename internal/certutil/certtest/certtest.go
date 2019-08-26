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

// Package certtest provides utilities for testing with certificate chains.
package certtest

import (
	"crypto"

	"github.com/WICG/webpackage/go/signedexchange/certurl"
	"github.com/google/webpackager/internal/certutil"
)

// ReadCertChainFile is like certutil.ReadCertChainFile but panics on error
// for ease of use in testing.
func ReadCertChainFile(filename string) certurl.CertChain {
	cert, err := certutil.ReadCertChainFile(filename)
	if err != nil {
		panic(err)
	}
	return cert
}

// ReadPrivateKeyFile is like certutil.ReadPrivateKeyFile but panics on error
// for ease of use in testing.
func ReadPrivateKeyFile(filename string) crypto.PrivateKey {
	key, err := certutil.ReadPrivateKeyFile(filename)
	if err != nil {
		panic(err)
	}
	return key
}
