// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tomlconfig

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/webpackager/internal/urlutil"
	"golang.org/x/xerrors"
)

const (
	// maxExpiry is the maximum period for SXG.Expiry.
	maxExpiry = 7 * (24 * time.Hour)

	// maxJSExpiry is the maximum period for SXG.JSExpiry. This limit can
	// be bypassed (up to maxExpiry) by adding "unsafe:" prefix.
	maxJSExpiry = 24 * time.Hour
)

// GetExpiry returns a parsed c.Expiry. It panics if c.Expiry contains an
// invalid value; it should not happen if c is obtained using ParseConfig
// or ReadFromFile.
func (c *SXGConfig) GetExpiry() time.Duration {
	d, err := parseExpiry(c.Expiry)
	if err != nil {
		panic(err)
	}
	return d
}

// GetJSExpiry returns a parsed c.JSExpiry. It panics if c.JSExpiry contains
// an invalid value; it should not happen if c is obtained using ParseConfig
// or ReadFromFile.
func (c *SXGConfig) GetJSExpiry() time.Duration {
	d, err := parseJSExpiry(c.JSExpiry)
	if err != nil {
		panic(err)
	}
	return d
}

func parseExpiry(value string) (time.Duration, error) {
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, errors.New("SXGs should have a positive lifetime")
	}
	if d > maxExpiry {
		maxHours := maxExpiry.Hours()
		return 0, xerrors.Errorf("SXGs must expire within %v hours", maxHours)
	}
	return d, nil
}

func parseJSExpiry(value string) (time.Duration, error) {
	trimmed := strings.TrimPrefix(value, "unsafe:")

	d, err := parseExpiry(trimmed)
	if err != nil {
		return 0, err
	}
	if d > maxJSExpiry && len(trimmed) == len(value) {
		return 0, xerrors.Errorf(`%q requires "unsafe:" prefix`, value)
	}
	return d, nil
}

// GetCertURLBase returns a parsed c.CertURLBase. It panics if c.CertURLBase
// cannot be parsed; it should not happen if c is obtained using ParseConfig
// or ReadFromFile.
func (c *SXGConfig) GetCertURLBase() *url.URL {
	return urlutil.MustParse(c.CertURLBase)
}

// GetValidityURL returns a parsed c.ValidityURL. It panics if c.ValidityURL
// cannot be parsed; it should not happen if c is obtained using ParseConfig
// or ReadFromFile.
func (c *SXGConfig) GetValidityURL() *url.URL {
	return urlutil.MustParse(c.ValidityURL)
}

// GetPathRE returns a compiled c.PathRE. It also encloses the regexp with
// `\A(?:...)\z` to make it a full match. It panics if c.PathRE is malformed;
// it should not happen if c is obtained using ParseConfig or ReadFromFile.
func (c *URLConfig) GetPathRE() *regexp.Regexp {
	return mustCompileFullMatch(c.PathRE)
}

// GetQueryRE returns a compiled c.QueryRE. It also encloses the regexp with
// `\A(?:...)\z` to make it a full match. It panics if c.PathRE is malformed;
// it should not happen if c is obtained using ParseConfig or ReadFromFile.
func (c *URLConfig) GetQueryRE() *regexp.Regexp {
	return mustCompileFullMatch(c.QueryRE)
}

func mustCompileFullMatch(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`\A(?:` + pattern + `)\z`)
}
