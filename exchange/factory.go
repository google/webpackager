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

package exchange

import (
	"crypto"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/WICG/webpackage/go/signedexchange/certurl"
	"github.com/WICG/webpackage/go/signedexchange/version"
	"github.com/google/webpackager/internal/certutil"
)

// Factory holds the parameters to generate signed exchanges.
type Factory struct {
	Version      version.Version
	MIRecordSize int
	CertChain    certurl.CertChain
	CertURL      *url.URL
	PrivateKey   crypto.PrivateKey
}

// NewExchange generates a signed exchange from resp, vp, and validityURL.
func (fty *Factory) NewExchange(resp *Response, vp ValidPeriod, validityURL *url.URL) (*signedexchange.Exchange, error) {
	e := signedexchange.NewExchange(
		fty.Version,
		resp.Request.URL.String(),
		resp.Request.Method,
		resp.Request.Header,
		resp.StatusCode,
		resp.GetFullHeader(),
		resp.Payload)
	if err := e.MiEncodePayload(fty.MIRecordSize); err != nil {
		return nil, err
	}

	signer := &signedexchange.Signer{
		Date:        vp.Date(),
		Expires:     vp.Expires(),
		Certs:       certutil.GetCertificates(fty.CertChain),
		CertUrl:     fty.CertURL,
		ValidityUrl: validityURL,
		PrivKey:     fty.PrivateKey,
	}
	if err := e.AddSignatureHeader(signer); err != nil {
		return nil, err
	}

	return e, nil
}

// Verify validates the provided signed exchange e at the provided date.
// It returns the payload decoded from e on success.
func (fty *Factory) Verify(e *signedexchange.Exchange, date time.Time) ([]byte, error) {
	var logText strings.Builder

	payload, ok := e.Verify(
		date,
		certutil.FakeCertFetcher(fty.CertChain, fty.CertURL),
		log.New(&logText, "", 0))
	if !ok {
		return nil, errors.New(logText.String())
	}

	return payload, nil
}
