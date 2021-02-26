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
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/crypto/ocsp"
)

// MaxOCSPResponseDuration represents the maximum duration allowed for the
// validity period of OCSP responses. used with signed exchanges.
const MaxOCSPResponseDuration = 7 * (24 * time.Hour) // 7 days

// ErrDummyOCSPResponse is returned if VerifyForRawChain/VerifySXGCriteria
// is called on DummyOCSPResponse.
var ErrDummyOCSPResponse = errors.New("certchain: verifying dummy OCSPResponse")

// OCSPResponse wraps an ocsp.Response with the DER bytes.
type OCSPResponse struct {
	*ocsp.Response
	Raw []byte
}

// DummyOCSPResponse is a dummy OCSPResponse to use with test certificates
// lacking OCSP responders, such as self-signed certificates.
//
// Note DummyOCSPResponse does not comprise a valid OCSP response. It just
// provides dummy bytes to fill in the application/cert-chain+cbor stream.
var DummyOCSPResponse *OCSPResponse = &OCSPResponse{
	new(ocsp.Response),
	[]byte("dummy-ocsp"),
}

// ParseOCSPResponse parses an OCSP response in DER form. It only supports
// responses for a single certificate. If the response contains a certificate
// then the signature over the response is checked.
func ParseOCSPResponse(bytes []byte) (*OCSPResponse, error) {
	resp, err := ocsp.ParseResponse(bytes, nil)
	if err != nil {
		if string(bytes) == string(DummyOCSPResponse.Raw) {
			// This is necessary to serve a(n invalid) cert-chain+cbor in
			// AllowTestCert mode.
			return DummyOCSPResponse, nil
		}
		return nil, err
	}
	return &OCSPResponse{resp, bytes}, nil
}

// ParseOCSPResponseForRawChain parses an OCSP response in DER form and
// searches for an OCSPResponse relating to c. If such an OCSPResponse is
// found and the OCSP response contains a certificate then the signature over
// the response is checked. c.Issuer will be used to validate the signature
// or embedded certificate.
func ParseOCSPResponseForRawChain(derBytes []byte, c *RawChain) (*OCSPResponse, error) {
	resp, err := ocsp.ParseResponseForCert(derBytes, c.Leaf, c.Issuer)
	if err != nil {
		return nil, err
	}
	return &OCSPResponse{resp, derBytes}, nil
}

// VerifyForRawChain verifies that resp is valid at the provided time t for
// the RawChain c. More specifically it checks resp has:
//
//   - a serial number matching c.Leaf.
//   - a valid signature or embedded certificate from c.Issuer.
//   - an update period that includes t.
//
// VerifyForRawChain returns ErrDummyOCSPResponse if resp is DummyOCSPResponse.
// In other error cases, VerifyForRawChain returns a multierror.Error
// (hashicorp/go-multierror) to report as many problems as possible.
//
// BUG(yuizumi): VerifyForRawChain should verify the OCSPResponse has both
// a matching serial number and a matching issuer, but it verifies the issuer
// only indirectly, through the signature or embedded certificate.
func (resp *OCSPResponse) VerifyForRawChain(t time.Time, c *RawChain) error {
	if resp == DummyOCSPResponse {
		return ErrDummyOCSPResponse
	}

	var errs *multierror.Error

	errs = multierror.Append(errs, verifyWithRawChain(resp, c))
	errs = multierror.Append(errs, verifyWithTime(resp, t))

	return errs.ErrorOrNil()
}

func verifyWithRawChain(resp *OCSPResponse, c *RawChain) error {
	if resp.SerialNumber.Cmp(c.Leaf.SerialNumber) != 0 {
		// ParseResponseForCert is meaningless this case.
		return errors.New("certchain: SerialNumber does not match")
	}
	_, err := ocsp.ParseResponseForCert(resp.Raw, c.Leaf, c.Issuer)
	return err
}

func verifyWithTime(resp *OCSPResponse, t time.Time) error {
	start := resp.ThisUpdate
	end := resp.NextUpdate

	// This is technically redundant but gives a finer error message.
	if !start.Before(end) {
		return fmt.Errorf("certchain: OCSP thisUpdate and nextUpdate are corrupted: %+v to %+v", start, end)
	}
	if t.Before(start) {
		return fmt.Errorf("certchain: OCSP thisUpdate is in the future: %+v", start)
	}
	if t.After(end) {
		return fmt.Errorf("certchain: OCSP nextUpdate is in the past: %+v", end)
	}

	return nil
}

// VerifySXGCriteria verifies that resp satisfies the criteria for use with
// signed exchanges. More specifically it checks resp has:
//
//   - ocsp.Good as its Status value.
//   - an update interval not longer than MaxOCSPResponseDuration.
//
// VerifySXGCriteria returns ErrDummyOCSPResponse if resp is DummyOCSPResponse.
// In other error cases, VerifySXGCriteria returns a multierror.Error
// (hashicorp/go-multierror) to report as many problems as possible.
func (resp *OCSPResponse) VerifySXGCriteria() error {
	if resp == DummyOCSPResponse {
		return ErrDummyOCSPResponse
	}

	var errs *multierror.Error

	if resp.Status != ocsp.Good {
		err := fmt.Errorf("certchain: invalid OCSP status: %+v", resp.Status)
		errs = multierror.Append(errs, err)
	}

	start := resp.ThisUpdate
	end := resp.NextUpdate
	if end.Sub(start) > MaxOCSPResponseDuration {
		err := fmt.Errorf("certchain: OCSP update interval too long: %+v to %+v", start, end)
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}
