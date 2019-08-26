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

// Package certutil provides utility functions to handle certificate chains.
package certutil

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/WICG/webpackage/go/signedexchange/certurl"
)

var (
	errUnknownCertURL = errors.New("certuril: unknown cert-url")
)

// FakeCertFetcher returns a CertFetcher that returns the serialized certChain
// when the url argument is equal to certURL and fails with an error otherwise.
func FakeCertFetcher(certChain certurl.CertChain, certURL *url.URL) signedexchange.CertFetcher {
	return func(url string) ([]byte, error) {
		if url != certURL.String() {
			return nil, errUnknownCertURL
		}
		var cborBuf bytes.Buffer
		if err := certChain.Write(&cborBuf); err != nil {
			return nil, err
		}
		return cborBuf.Bytes(), nil
	}
}

// FetchCertChain downloads a CertChain from the given URL.
func FetchCertChain(url *url.URL) (certurl.CertChain, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return certurl.ReadCertChain(resp.Body)
}

// GetCertificates retrieves all x509.Certificate instances from a CertChain.
func GetCertificates(certChain certurl.CertChain) []*x509.Certificate {
	certificates := make([]*x509.Certificate, len(certChain))
	for i, item := range certChain {
		certificates[i] = item.Cert
	}
	return certificates
}

// ReadCertChainFile reads a certificate chain from a file in the CBOR format
// (application/cert-chain+cbor).
func ReadCertChainFile(filename string) (certurl.CertChain, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return certurl.ReadCertChain(f)
}

// ReadCertificatesFile reads a series of x509.Certificate from a PEM file.
func ReadCertificatesFile(filename string) ([]*x509.Certificate, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return signedexchange.ParseCertificates(content)
}

// ReadPrivateKeyFile reads a PEM file and returns a PrivateKey usable to sign
// exchanges. ReadPrivateKeyFile returns an error when the PEM file contains no
// usable key.
func ReadPrivateKeyFile(filename string) (crypto.PrivateKey, error) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return signedexchange.ParsePrivateKey(text)
}
