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

// Package certchaintest provides utilities for certificate chain testing.
package certchaintest

import (
	"crypto"

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certchainutil"
)

// MustReadRawChainFile is like certchainutil.ReadRawChainFile but panics
// on error for ease of use in testing.
func MustReadRawChainFile(filename string) *certchain.RawChain {
	got, err := certchainutil.ReadRawChainFile(filename)
	if err != nil {
		panic(err)
	}
	return got
}

// MustReadOCSPRespFile is like certchainutil.ReadOCSPRespFile but panics
// on error for ease of use in testing.
func MustReadOCSPRespFile(filename string) *certchain.OCSPResponse {
	got, err := certchainutil.ReadOCSPRespFile(filename)
	if err != nil {
		panic(err)
	}
	return got
}

// MustReadAugmentedChainFile is like certchainutil.ReadAugmentedChainFile
// but panics on error for ease of use in testing.
func MustReadAugmentedChainFile(filename string) *certchain.AugmentedChain {
	got, err := certchainutil.ReadAugmentedChainFile(filename)
	if err != nil {
		panic(err)
	}
	return got
}

// MustReadPrivateKeyFile is like certchainutil.ReadPrivateKeyFile but panics
// on error for ease of use in testing.
func MustReadPrivateKeyFile(filename string) crypto.PrivateKey {
	got, err := certchainutil.ReadPrivateKeyFile(filename)
	if err != nil {
		panic(err)
	}
	return got
}
