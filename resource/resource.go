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

// Package resource defines representations of resources to generate signed
// exchanges for.
package resource

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/WICG/webpackage/go/signedexchange"
)

// Resource represents a resource for which a signed exchange is generated.
type Resource struct {
	// RequestURL is the location to request this resource.
	RequestURL *url.URL

	// RedirectURL is non-nil when RequestURL gets redirected, and the value
	// represents the location after the redirect.
	//
	// All other fields, except RequestURL, should be unset when RedirectURL
	// has a non-nil value.
	RedirectURL *url.URL

	// PhysicalURL is the physical location of this resource. It is usually
	// identical to RequestURL, but can be different when the web server has
	// URL rewriting mechanism (e.g. to serve "index.html").
	//
	// See also: Package urlrewrite.
	PhysicalURL *url.URL

	// ValidityURL represents the location of the validity data.
	ValidityURL *url.URL

	// Exchange represents a signed exchange generated for this resource.
	//
	// Exchange should be set through the SetExchange method to keep in sync
	// with the Integrity field.
	Exchange *signedexchange.Exchange

	// Integrity represents the integirty of HTTP response headers for this
	// resource. Technically, it is the hash of the CBOR representation of
	// the response headers in the signed exchange, prefixed by the algorithm
	// and base64-encoded ("sha256-wZp5f6H9...").
	//
	// Integrity is set by the SetExchange method.
	Integrity string
}

// NewResource creates and initializes a new Resource for url.
func NewResource(url *url.URL) *Resource {
	return &Resource{RequestURL: url}
}

// String returns a string representing the Resource.
func (r *Resource) String() string {
	return fmt.Sprintf("<%s>", r.RequestURL)
}

// SetExchange populates r.Exchange with the provided signed exchange e and
// updates r.Integrity accordingly. It returns a non-nil error when it fails
// to compute the new integrity value, in which case r is not mutated.
func (r *Resource) SetExchange(e *signedexchange.Exchange) error {
	hasher := sha256.New()

	if err := e.DumpExchangeHeaders(hasher); err != nil {
		return err
	}
	sum := hasher.Sum(nil)

	r.Exchange = e
	r.Integrity = "sha256-" + base64.StdEncoding.EncodeToString(sum)
	return nil
}

// AllowedAltSXGHeader returns the value of a Link HTTP header to allow the
// resource to be distributed from different domains (rel="allowed-alt-sxg").
func (r *Resource) AllowedAltSXGHeader() string {
	// TODO(yuizumi): Support Variants.
	return fmt.Sprintf(`<%s>;rel="allowed-alt-sxg";header-integrity=%q`,
		r.RequestURL, r.Integrity)
}
