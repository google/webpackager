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

// Package certchainutil complements the certchain package.
package certchainutil

import (
	"bytes"
	"crypto"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/google/webpackager/certchain"
)

// FetchAugmentedChain retrieves an AugmentedChain from url.
func FetchAugmentedChain(url *url.URL) (*certchain.AugmentedChain, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return certchain.ReadAugmentedChain(resp.Body)
}

// ReadRawChainFile reads a PEM file to retrieve a RawChain.
func ReadRawChainFile(filename string) (*certchain.RawChain, error) {
	pem, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return certchain.NewRawChainFromPEM(pem)
}

// ReadOCSPRespFile reads a DER file to retrieve an OCSPResponse, using
// certchain.ParseOCSPResponse.
func ReadOCSPRespFile(filename string) (*certchain.OCSPResponse, error) {
	der, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return certchain.ParseOCSPResponse(der)
}

// ReadAugmentedChainFile reads an application/cert-chain+cbor file to
// retrieve an AugmentedChain.
func ReadAugmentedChainFile(filename string) (*certchain.AugmentedChain, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return certchain.ReadAugmentedChain(f)
}

// ReadPrivateKeyFile reads a PEM file and returns a PrivateKey usable
// for signing exchanges.
func ReadPrivateKeyFile(filename string) (crypto.PrivateKey, error) {
	pem, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return signedexchange.ParsePrivateKey(pem)
}

// WrapToCertFetcher wraps an AugmentedChain into a signedexchange.CertFetcher.
// The CertFetcher does not inspect the url argument.
func WrapToCertFetcher(c *certchain.AugmentedChain) signedexchange.CertFetcher {
	return func(url string) ([]byte, error) {
		var b bytes.Buffer
		if err := c.WriteCBOR(&b); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}
}
