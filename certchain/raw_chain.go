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

package certchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/hashicorp/go-multierror"
)

// MaxCertDuration represents the maximum duration allowed for the validity
// period of signed exchange certificates.
const MaxCertDuration = 90 * (24 * time.Hour) // 90 days

// RawChain represents an X509 certificate chain, populated with information
// extracted from it for convenience.
type RawChain struct {
	// Certs is the array of certificates which form this certificate chain,
	// starting with the end-entity certificate.
	Certs []*x509.Certificate

	// Digest gives a unique identifier of this certificate chain, produced
	// using a hash function.
	Digest string

	// Leaf represents the end-entity certificate of this certificate chain.
	// It is always equal to Certs[0].
	Leaf *x509.Certificate

	// Issuer represents the certificate of the Leaf's direct issuer. It is
	// equal to Certs[1] for CA-issued certificates and Certs[0] (Leaf) for
	// self-signed certificates.
	Issuer *x509.Certificate

	// OCSPServer is the URI of the Leaf's OCSP responder. If Leaf does not
	// have an OCSP responder, OCSPServer is an empty string.
	OCSPServer string
}

// NewRawChain creates a new RawChain with certs.
//
// certs must form a certificate chain, where the first element is the
// end-entity certificate and the last element is the root certificate
// or the certificate issued by a trusted root. Each certificate in the
// chain must be followed by the certificate of its direct issuer, except
// for the last certificate.
func NewRawChain(certs []*x509.Certificate) (*RawChain, error) {
	if len(certs) == 0 {
		return nil, errors.New("certchain: empty certs")
	}

	leaf := certs[0]

	for i, cert := range certs[:len(certs)-1] {
		issuer := certs[i+1]

		// BUG(yuizumi): We are using bytes.Equal to match the issuer and
		// the subject, like the crypto/x509 package. It is not the way
		// we are supposed to compare distinguished names, although it is
		// a good approximate.
		if !bytes.Equal(cert.RawIssuer, issuer.RawSubject) {
			return nil, errors.New("certchain: certs not forming a chain")
		}
	}

	var issuer *x509.Certificate
	var ocspServer string

	if len(certs) >= 2 {
		issuer = certs[1]
	} else {
		issuer = certs[0]
	}
	if len(leaf.OCSPServer) >= 1 {
		ocspServer = leaf.OCSPServer[0]
	}

	c := &RawChain{certs, digest(certs), leaf, issuer, ocspServer}
	return c, nil
}

// NewRawChainFromPEM creates a new RawChain from PEM bytes.
func NewRawChainFromPEM(bytes []byte) (*RawChain, error) {
	certs, err := signedexchange.ParseCertificates(bytes)
	if err != nil {
		return nil, err
	}
	return NewRawChain(certs)
}

func digest(certs []*x509.Certificate) string {
	hasher := sha256.New()
	for _, cert := range certs {
		hasher.Write(cert.Raw)
	}
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
}

// WritePEM writes the RawChain to w in the PEM format.
func (c *RawChain) WritePEM(w io.Writer) error {
	for _, cert := range c.Certs {
		block := &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
		if err := pem.Encode(w, block); err != nil {
			return err
		}
	}
	return nil
}

// VerifyChain attempts to verify that c is valid as of the provided time t,
// calling c.Leaf.Verify internally.
//
// WARNING: VerifyChain does not verify that the root certificate is trusted
// by operating systems or user agents.
func (c *RawChain) VerifyChain(t time.Time) error {
	opts := x509.VerifyOptions{
		Roots:         x509.NewCertPool(),
		Intermediates: x509.NewCertPool(),
		CurrentTime:   t,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	for i, cert := range c.Certs {
		// Treat the last cert always as a "root" to terminate the chain.
		if i == len(c.Certs)-1 {
			opts.Roots.AddCert(cert)
		} else {
			opts.Intermediates.AddCert(cert)
		}
	}
	_, err := c.Leaf.Verify(opts)

	return err
}

// VerifySXGCriteria verifies that the RawChain c satisifes the criteria
// for use with signed exchanges. More specifically it checks c.Leaf has:
//
//   - a public key of supported cryptographic algorithm.
//   - canHttpSignExchange extension.
//   - a validity period not longer than MaxCertDuration.
//
// VerifySXGCriteria returns multierror.Error (hashicorp/go-multierror) to
// report as many problems as possible.
//
// BUG(yuizumi): VerifySXGCriteria accepts only ECDSA-P256 and ECDSA-P384
// public keys; the signedexchange package (WICG/webpackage) supports only
// those keys at the moment.
func (c *RawChain) VerifySXGCriteria() error {
	var errs *multierror.Error

	errs = multierror.Append(errs, verifyPublicKey(c.Leaf))
	errs = multierror.Append(errs, verifyExtensions(c.Leaf))
	errs = multierror.Append(errs, verifyCertDuration(c.Leaf))

	return errs.ErrorOrNil()
}

func verifyPublicKey(cert *x509.Certificate) error {
	algo := cert.PublicKeyAlgorithm

	switch algo {
	case x509.RSA:
		return fmt.Errorf("certchain: invalid PublicKeyAlgorithm: %s", algo)
	case x509.DSA:
		return fmt.Errorf("certchain: invalid PublicKeyAlgorithm: %s", algo)
	case x509.ECDSA:
		key := cert.PublicKey.(*ecdsa.PublicKey)
		switch key.Curve {
		case elliptic.P256():
			return nil
		case elliptic.P384():
			return nil
		default:
			return fmt.Errorf("certchain: unknown elliptic curve: %s", key.Curve.Params().Name)
		}
	default:
		return fmt.Errorf("certchain: unknown PublicKeyAlgorithm: %s", algo)
	}
}

func verifyExtensions(cert *x509.Certificate) error {
	if findExtension(cert.Extensions, oidCanSignHttpExchanges) == nil {
		return errors.New("certchain: missing canSignHttpExchanges")
	}
	return nil
}

func verifyCertDuration(cert *x509.Certificate) error {
	s := cert.NotBefore
	t := cert.NotAfter
	if t.Sub(s) <= 0 {
		return fmt.Errorf("certchain: validity period corrupted: %+v to %+v", s, t)
	}
	if t.Sub(s) > MaxCertDuration {
		return fmt.Errorf("certchain: validity period too long: %+v to %+v", s, t)
	}
	return nil
}
