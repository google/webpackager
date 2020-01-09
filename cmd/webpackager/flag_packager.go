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

	"github.com/WICG/webpackage/go/signedexchange/version"
	"github.com/google/webpackager"
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/internal/certutil"
	"github.com/google/webpackager/internal/customflag"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/complexproc"
	"github.com/google/webpackager/resource/cache"
	"github.com/google/webpackager/resource/cache/filewrite"
	"github.com/google/webpackager/urlrewrite"
	"github.com/google/webpackager/validity"
	multierror "github.com/hashicorp/go-multierror"
)

var (
	// RequestTweaker
	flagRequestHeader = customflag.MultiString("request_header", `Request headers, e.g. "Accept-Language: en-US, en;q=0.5". (repeatable)`)

	// ExchangeFactory
	flagVersion      = flag.String("version", "1b3", `Signed exchange version.`)
	flagMIRecordSize = flag.String("mi_record_size", "4096", `Merkle Integration content encoding record size.`)
	flagCertCBOR     = flag.String("cert_cbor", "", `Certificate chain CBOR file. Fetched from --cert_url when unspecified.`)
	flagCertURL      = flag.String("cert_url", "", `Certficiate chain URL. (required)`)
	flagPrivateKey   = flag.String("private_key", "", `Private key PEM file. (required)`)

	// Processor
	flagSizeLimit = flag.String("size_limit", "4194304", fmt.Sprintf(`Maximum size of resources in bytes allowed for signed exchanges, or %q to set no limit.`, noSizeLimitString))

	// PhysicalURLRule
	flagIndexFile = flag.String("index_file", "index.html", `Filename assumed for slash-ended URLs.`)

	// ResourceCache, ValidityURLRule
	flagSXGExt      = flag.String("sxg_ext", ".sxg", `File extension for signed exchange files.`)
	flagSXGDir      = flag.String("sxg_dir", "sxg/", `Directory to output signed exchange files.`)
	flagValidityExt = flag.String("validity_ext", ".validity", `File extension for validity files. Note it is followed by a UNIX timestamp.`)
	flagValidityDir = flag.String("validity_dir", "", `Directory to output validity files. (unimplemented)`)
)

const (
	noSizeLimitString = "none"
)

func getConfigFromFlags() (*webpackager.Config, error) {
	cfg := new(webpackager.Config)
	var err error
	errs := new(multierror.Error)

	cfg.RequestTweaker, err = getRequestTweakerFromFlags()
	errs = multierror.Append(errs, err)
	cfg.PhysicalURLRule, err = getPhysicalURLRuleFromFlags()
	errs = multierror.Append(errs, err)
	cfg.ValidityURLRule, err = getValidityURLRuleFromFlags()
	errs = multierror.Append(errs, err)
	cfg.Processor, err = getProcessorFromFlags()
	errs = multierror.Append(errs, err)
	cfg.ExchangeFactory, err = getExchangeFactoryFromFlags()
	errs = multierror.Append(errs, err)
	cfg.ResourceCache, err = getResourceCacheFromFlags()
	errs = multierror.Append(errs, err)

	if err := errs.ErrorOrNil(); err != nil {
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

func parseByteSize(s string) (int, error) {
	// TODO(yuizumi): Maybe support binary suffixes (e.g. "4k" == 4096).
	v, err := strconv.Atoi(s)
	if err != nil {
		return v, err
	}
	if v <= 0 {
		return v, errors.New("value must be positive")
	}
	return v, nil
}

func parseSizeLimit(s string) (int, error) {
	if s == noSizeLimitString {
		return -1, nil
	}
	return parseByteSize(s)
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

func getRequestTweakerFromFlags() (fetch.RequestTweaker, error) {
	header := make(http.Header)
	errs := new(multierror.Error)

	for _, s := range *flagRequestHeader {
		chunks := strings.SplitN(s, ":", 2)
		if len(chunks) == 2 {
			key := strings.TrimSpace(chunks[0])
			val := strings.TrimSpace(chunks[1])
			header.Add(key, val)
		} else {
			errs = multierror.Append(
				errs, fmt.Errorf("invalid --request_header %q", s))
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	t := fetch.DefaultRequestTweaker
	if len(header) != 0 {
		t = fetch.RequestTweakerSequence{t, fetch.SetCustomHeaders(header)}
	}
	return t, nil
}

func getPhysicalURLRuleFromFlags() (urlrewrite.Rule, error) {
	rule := urlrewrite.RuleSequence{
		urlrewrite.CleanPath(),
		urlrewrite.IndexRule(*flagIndexFile),
	}
	return rule, nil
}

func getValidityURLRuleFromFlags() (validity.ValidityURLRule, error) {
	return validity.AppendExtDotLastModified(*flagValidityExt), nil
}

func getProcessorFromFlags() (processor.Processor, error) {
	var cfg complexproc.Config
	var err error
	errs := new(multierror.Error)

	cfg.Preverify.MaxContentLength, err = parseSizeLimit(*flagSizeLimit)
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("invalid --size_limit: %v", err))
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}
	return complexproc.NewComprehensiveProcessor(cfg), nil
}
func getExchangeFactoryFromFlags() (*exchange.Factory, error) {
	fty := new(exchange.Factory)
	var err error
	errs := new(multierror.Error)

	fty.Version, err = parseVersion(*flagVersion)
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("invalid --version: %v", err))
	}

	fty.MIRecordSize, err = parseByteSize(*flagMIRecordSize)
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("invalid --mi_record_size: %v", err))
	}

	if *flagCertURL == "" {
		errs = multierror.Append(errs, errors.New("missing --cert_url"))
	} else {
		fty.CertURL, err = parseCertURL(*flagCertURL)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("invalid --cert_url: %v", err))
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
		errs = multierror.Append(errs, fmt.Errorf("failed to load cert chain from %q: %v", certChainSource, err))
	}

	if *flagPrivateKey == "" {
		errs = multierror.Append(errs, errors.New("missing --private_key"))
	} else {
		fty.PrivateKey, err = certutil.ReadPrivateKeyFile(*flagPrivateKey)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to load private key from %q: %v", *flagPrivateKey, err))
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}
	return fty, err
}

func getResourceCacheFromFlags() (cache.ResourceCache, error) {
	config := filewrite.Config{BaseCache: cache.NewOnMemoryCache()}

	if *flagSXGDir != "" {
		config.ExchangeMapping = filewrite.AddBaseDir(
			filewrite.AppendExt(filewrite.UsePhysicalURLPath(), *flagSXGExt),
			*flagSXGDir,
		)
	}
	if *flagValidityDir != "" {
		return nil, errors.New("--validity_dir is not implemented yet")
	}

	return filewrite.NewFileWriteCache(config), nil
}
