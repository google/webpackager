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

package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/WICG/webpackage/go/signedexchange/version"
	"github.com/google/webpackager"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/internal/certutil"
	"github.com/google/webpackager/internal/customflag"
	"github.com/google/webpackager/internal/multierror"
	"github.com/google/webpackager/resource/cache"
	"github.com/google/webpackager/resource/cache/filewrite"
	"github.com/google/webpackager/urlrewrite"
	"github.com/google/webpackager/validity"
)

var (
	// RequestHeader
	flagRequestHeader = customflag.MultiString("request_header", `Request headers, e.g. "Accept-Language: en-US, en;q=0.5". (repeatable)`)

	// ExchangeFactory
	flagVersion      = flag.String("version", "1b3", `Signed exchange version.`)
	flagMIRecordSize = flag.String("mi_record_size", "4096", `Merkle Integration content encoding record size.`)
	flagCertCBOR     = flag.String("cert_cbor", "", `Certificate chain CBOR file. Fetched from --cert_url when unspecified.`)
	flagCertURL      = flag.String("cert_url", "", `Certficiate chain URL. (required)`)
	flagPrivateKey   = flag.String("private_key", "", `Private key PEM file. (required)`)

	// PhysicalURLRule
	flagIndexFile = flag.String("index_file", "index.html", `Filename assumed for slash-ended URLs.`)

	// ResourceCache, ValidityURLRule
	flagSXGExt      = flag.String("sxg_ext", ".sxg", `File extension for signed exchange files.`)
	flagSXGDir      = flag.String("sxg_dir", "sxg/", `Directory to output signed exchange files.`)
	flagValidityExt = flag.String("validity_ext", ".validity", `File extension for validity files. Note it is followed by a UNIX timestamp.`)
	flagValidityDir = flag.String("validity_dir", "", `Directory to output validity files. (unimplemented)`)
)

const (
	dateNowString = "now"
	maxExpiry     = 7 * (24 * time.Hour)
)

func getConfigFromFlags() (*webpackager.Config, error) {
	var errs multierror.MultiError
	cfg := &webpackager.Config{
		RequestHeader:   getRequestHeaderFromFlags(&errs),
		PhysicalURLRule: getPhysicalURLRuleFromFlags(&errs),
		ValidityURLRule: getValidityURLRuleFromFlags(&errs),
		ExchangeFactory: getExchangeFactoryFromFlags(&errs),
		ResourceCache:   getResourceCacheFromFlags(&errs),
	}
	if err := errs.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseVersion(s string) (version.Version, error) {
	v, ok := version.Parse(s)
	if !ok {
		return "", errors.New("unknown version")
	}
	return v, nil
}

func parseMIRecordSize(s string) (int, error) {
	// TODO(yuizumi): Maybe support binary prefixes (e.g. "4k" == 4096).
	v, err := strconv.Atoi(s)
	if err != nil {
		return v, err
	}
	if v <= 0 {
		return v, errors.New("value must be positive")
	}
	return v, nil
}

func parseCertURL(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, errors.New("must be an https:// url")
	}
	return u, nil
}

func getRequestHeaderFromFlags(errs *multierror.MultiError) http.Header {
	header := http.Header{}

	for _, s := range *flagRequestHeader {
		chunks := strings.SplitN(s, ":", 2)
		if len(chunks) == 2 {
			key := strings.TrimSpace(chunks[0])
			val := strings.TrimSpace(chunks[1])
			header.Add(key, val)
		} else {
			errs.Add(fmt.Errorf("invalid --request_header %q", s))
		}
	}

	return header
}

func getPhysicalURLRuleFromFlags(errs *multierror.MultiError) urlrewrite.Rule {
	return urlrewrite.RuleSequence{
		urlrewrite.CleanPath(),
		urlrewrite.IndexRule(*flagIndexFile),
	}
}

func getValidityURLRuleFromFlags(errs *multierror.MultiError) validity.ValidityURLRule {
	// err will be logged in getExchangeFactoryFromFlags().
	date, err := parseDate(*flagDate)
	if err != nil {
		date = time.Now()
	}
	return validity.AppendExtDotUnixTime(*flagValidityExt, date)
}

func getExchangeFactoryFromFlags(errs *multierror.MultiError) *exchange.Factory {
	fty := &exchange.Factory{}
	var err error

	fty.Version, err = parseVersion(*flagVersion)
	if err != nil {
		errs.Add(fmt.Errorf("invalid --version: %v", err))
	}

	fty.MIRecordSize, err = parseMIRecordSize(*flagMIRecordSize)
	if err != nil {
		errs.Add(fmt.Errorf("invalid --mi_record_size: %v", err))
	}

	if *flagCertURL == "" {
		errs.Add(errors.New("missing --cert_url"))
	} else {
		fty.CertURL, err = parseCertURL(*flagCertURL)
		if err != nil {
			errs.Add(fmt.Errorf("invalid --cert_url: %v", err))
		}
	}

	var certChainSource string
	if *flagCertCBOR != "" {
		fty.CertChain, err = certutil.ReadCertChainFile(*flagCertCBOR)
		certChainSource = *flagCertCBOR
	} else if fty.CertURL != nil {
		fty.CertChain, err = certutil.FetchCertChain(fty.CertURL)
		certChainSource = fty.CertURL.String()
	}
	if err != nil {
		errs.Add(fmt.Errorf("failed to load cert chain from %q: %v", certChainSource, err))
	}

	if *flagPrivateKey == "" {
		errs.Add(errors.New("missing --private_key"))
	} else {
		fty.PrivateKey, err = certutil.ReadPrivateKeyFile(*flagPrivateKey)
		if err != nil {
			errs.Add(fmt.Errorf("failed to load private key from %q: %v", *flagPrivateKey, err))
		}
	}

	return fty
}

func getResourceCacheFromFlags(errs *multierror.MultiError) cache.ResourceCache {
	config := filewrite.Config{BaseCache: cache.NewOnMemoryCache()}

	if *flagSXGDir != "" {
		config.ExchangeMapping = filewrite.AddBaseDir(
			filewrite.AppendExt(filewrite.UsePhysicalURLPath(), *flagSXGExt),
			*flagSXGDir,
		)
	}
	if *flagValidityDir != "" {
		errs.Add(errors.New("--validity_dir is not implemented yet"))
	}

	return filewrite.NewFileWriteCache(config)
}
