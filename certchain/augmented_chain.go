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
	"crypto/x509"
	"errors"
	"io"
	"time"

	"github.com/WICG/webpackage/go/signedexchange/certurl"
	"github.com/hashicorp/go-multierror"
)

// ErrInvalidOCSPValue is returned by ReadCBOR if the provided CBOR stream
// contained invalid OCSP response.
var ErrInvalidOCSPValue = errors.New("certchain: invalid ocsp value")

// AugmentedChain is a certificate chain augmented with an OCSP response
// and unembedded SCTs (Signed Certificate Timestamps) for the end-entity
// certificate. It is designed to support application/cert-chain+cbor
// certificate chains, but augments the certificate chain instead of each
// certificate. In particular, AugmentedChain stores unembedded SCTs only for
// the end-entity certificate while the application/cert-chain+cbor format
// can contain SCTs for every certificate. This difference should not matter
// in practice: the signed exchange validation process only uses SCTs of the
// end-entity certificate.
//
// AugmentedChain handles SCT lists as an opaque byte sequence. It does not
// know about the validity of SCTs against the certificate, for example.
type AugmentedChain struct {
	*RawChain

	// OCSPResp contains an OCSP response for the end-entity certificate.
	OCSPResp *OCSPResponse

	// SCTList contains unembedded SCTs for the end-entity certificate.
	//
	// Note SCTs can also be embedded in certificates and OCSP responses.
	// SCTList is required only when neither the end-entity certificate nor
	// its OCSP response contains embedded SCTs.
	SCTList []byte
}

// NewAugmentedChain creates a new AugmentedChain.
func NewAugmentedChain(c *RawChain, ocsp *OCSPResponse, sct []byte) *AugmentedChain {
	return &AugmentedChain{c, ocsp, sct}
}

// NewAugmentedChainFromCBOR creates a new AugmentedChain from a serialized
// certificate chain in the application/cert-chain+cbor format.
//
// If you are reading the certificate chain from a file or over the network,
// consider using ReadAugmentedChain. It stops reading immediately when it has
// detected an error in the middle.
//
// See ReadAugmentedChain for how the ocsp and sct values are handled.
func NewAugmentedChainFromCBOR(cborBytes []byte) (*AugmentedChain, error) {
	return ReadAugmentedChain(bytes.NewReader(cborBytes))
}

// ReadAugmentedChain reads an application/cert-chain+cbor stream from r to
// create an AugmentedChain.
//
// The ocsp value is parsed into an OCSPResponse. In case of a parse error,
// ReadAugmentedChain creates an AugmentedChain with DummyOCSPResponse and
// returns it with ErrInvalidOCSPValue; the invalid ocsp value is discarded.
// The caller may expect or ignore ErrInvalidOCSPValue, e.g. when using a test
// certificate.
//
// ReadAugmentedChain keeps the sct value only for the end-entity certificate.
// The sct values for other certificates, if any, are silently discarded. Note
// AugmentedChain stores unembedded SCTs only for the end-entity certifiacte.
func ReadAugmentedChain(r io.Reader) (*AugmentedChain, error) {
	items, err := certurl.ReadCertChain(r)
	if err != nil {
		return nil, err
	}
	certs := make([]*x509.Certificate, len(items))
	for i, item := range items {
		certs[i] = item.Cert
	}
	c, err := NewRawChain(certs)
	if err != nil {
		return nil, err
	}
	ocsp, err := ParseOCSPResponseForRawChain(items[0].OCSPResponse, c)
	if err != nil {
		ac := NewAugmentedChain(c, DummyOCSPResponse, items[0].SCTList)
		return ac, ErrInvalidOCSPValue
	}
	return NewAugmentedChain(c, ocsp, items[0].SCTList), nil
}

// HasSCTList reports whether the AugmentedChain ac contains SCTs. It looks
// for an SCT extension in the end-entity certificate and the OCSP response
// for embedded SCTs, as well as the SCTList field for unembedded SCTs.
//
// HasSCTList only checks the existence, not the content. The SCTList field
// is assumed to contain SCTs unless it is nil or empty.
func (ac *AugmentedChain) HasSCTList() bool {
	if findExtension(ac.Certs[0].Extensions, oidCertEmbeddedSCT) != nil {
		return true
	}
	if findExtension(ac.OCSPResp.Extensions, oidOCSPEmbeddedSCT) != nil {
		return true
	}
	return len(ac.SCTList) != 0
}

// WriteCBOR writes an AugmentedChain to w in the application/cert-chain+cbor
// format.
func (ac *AugmentedChain) WriteCBOR(w io.Writer) error {
	cc, err := certurl.NewCertChain(ac.Certs, ac.OCSPResp.Raw, ac.SCTList)
	if err != nil {
		return err
	}
	return cc.Write(w)
}

// VerifyAll does comprehensive checks with ac. More specifically it checks:
//
//   - ac.RawChain.VerifyChain succeeds.
//   - ac.RawChain.VerifySXGCriteria succeeds.
//   - ac.OCSPResp.VerifyForRawChain succeeds.
//   - ac.OCSPResp.VerifySXGCriteria succeeds.
//   - ac.HasSCTList returns true.
// If inProduction is true, allow test certs and OCSP to have dummy value.
//
// VerifyAll returns a multierror.Error (hashicorp/go-multierror) to report
// as many problems as possible.
func (ac *AugmentedChain) VerifyAll(t time.Time, inProduction bool) error {
	var errs *multierror.Error

	errs = multierror.Append(errs, ac.RawChain.VerifyChain(t))
	errs = multierror.Append(errs, ac.RawChain.VerifySXGCriteria())

	if inProduction {
		if ac.OCSPResp == DummyOCSPResponse {
			// Avoid ErrDummyOCSResponse being appended twice.
			errs = multierror.Append(errs, ErrDummyOCSPResponse)
		} else {
			errs = multierror.Append(errs, ac.OCSPResp.VerifyForRawChain(t, ac.RawChain))
			errs = multierror.Append(errs, ac.OCSPResp.VerifySXGCriteria())
		}

		if !ac.HasSCTList() {
			errs = multierror.Append(errs, errors.New("certchain: missing SCTs"))
		}
	}

	return errs.ErrorOrNil()
}
