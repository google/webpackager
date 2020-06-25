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
	"net/url"

	"github.com/WICG/webpackage/go/signedexchange/version"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/internal/urlutil"
)

// These are the default values in Config.
const (
	DefaultVersion      = version.Version1b3
	DefaultMIRecordSize = 16384
)

// DefaultCertURL is the default value for CertURL in Config.
var DefaultCertURL = urlutil.MustParse("/cert.cbor")

// Config holds the parameters to produce signed exchanges.
type Config struct {
	// Version specifies the signed exchange version. If Version is empty,
	// Factory uses DefaultVersion.
	Version version.Version

	// MIRecordSize specifies Merkle Integrity record size. The value must
	// be positive, or zero to use DefaultMIRecordSize. It must not exceed
	// 16384 (16 KiB) to be compliant with the specification.
	MIRecordSize int

	// CertChain specifies the certificate chain. CertChain may not be nil.
	CertChain *certchain.AugmentedChain

	// CertURL specifies the cert-url parameter in the signature. It can be
	// relative, in which case Factory resolves the absolute URL using
	// the request URL. It should still usually contain an absolute path
	// (e.g. "/cert.cbor", not "cert.cbor"). If CertURL is nil, Factory uses
	// DefaultCertURL.
	CertURL *url.URL

	// PrivateKey specifies the private key used for signing. PrivateKey may
	// not be nil.
	PrivateKey crypto.PrivateKey

	// KeepNonSXGPreloads instructs Factory to include preload link headers
	// that don't have the corresponding allowed-alt-sxg with a valid
	// header-integrity.
	KeepNonSXGPreloads bool
}

func (c *Config) populateDefaults() {
	if c.CertChain == nil {
		panic("c.CertChain is nil")
	}
	if c.PrivateKey == nil {
		panic("c.PrivateKey is nil")
	}

	if c.Version == "" {
		c.Version = DefaultVersion
	}
	if c.MIRecordSize == 0 {
		c.MIRecordSize = DefaultMIRecordSize
	}
	if c.CertURL == nil {
		c.CertURL = DefaultCertURL
	}
}
