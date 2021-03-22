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

package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/webpackager/certchain/certmanager/acmeclient"

	"github.com/google/webpackager"
	"github.com/google/webpackager/certchain/certchainutil"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/exchange/vprule"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/complexproc"
	"github.com/google/webpackager/processor/htmlproc"
	"github.com/google/webpackager/processor/htmlproc/htmltask"
	"github.com/google/webpackager/processor/preverify"
	"github.com/google/webpackager/resource/cache"
	"github.com/google/webpackager/server/tomlconfig"
	"github.com/google/webpackager/urlmatcher"
	"github.com/google/webpackager/validity"
	"github.com/hashicorp/go-multierror"
)

// FromTOMLConfig creates and initializes a Server from TOML config.
func FromTOMLConfig(c *tomlconfig.Config) (*Server, error) {
	var errs *multierror.Error

	tlsConfig, err := makeTLSConfig(c)
	errs = multierror.Append(errs, err)
	exchangeFactory, err := makeExchangeFactory(c)
	errs = multierror.Append(errs, err)

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: net.JoinHostPort(
			c.Listen.Host,
			strconv.Itoa(c.Listen.Port),
		),
		TLSConfig: tlsConfig,
		// TODO(yuizumi): Allow changing timeouts in config.
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	pc := webpackager.Config{
		FetchClient:     makeFetchClient(c),
		ValidityURLRule: makeValidityURLRule(c),
		Processor:       makeProcessor(c),
		ValidPeriodRule: makeValidPeriodRule(c),
		ExchangeFactory: exchangeFactory,
	}

	if size := c.Cache.MaxEntries; size > 0 {
		pc.ResourceCache = cache.NewBoundedCache(size)
	} else if size == 0 {
		pc.ResourceCache = cache.NilCache()
	} else {
		pc.ResourceCache = cache.NewOnMemoryCache() // unbounded
	}

	config := Config{
		Packager:     webpackager.NewPackager(pc),
		CertManager:  exchangeFactory.CertManager,
		ServerConfig: c.Server,
	}

	return NewServer(server, config), nil
}

func makeTLSConfig(c *tomlconfig.Config) (*tls.Config, error) {
	if c.Listen.TLS.PEMFile == "" && c.Listen.TLS.KeyFile == "" {
		return nil, nil
	}
	cert, err := tls.LoadX509KeyPair(
		c.Listen.TLS.PEMFile,
		c.Listen.TLS.KeyFile,
	)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func makeFetchClient(c *tomlconfig.Config) fetch.FetchClient {
	allow := make([]urlmatcher.Matcher, len(c.Sign))
	for i, uc := range c.Sign {
		allow[i] = urlmatcher.AllOf(
			urlmatcher.HasScheme("https"),
			urlmatcher.HasHostname(uc.Domain),
			urlmatcher.HasEscapedPathRegexp(uc.GetPathRE()),
			urlmatcher.HasRawQueryRegexp(uc.GetQueryRE()),
		)
	}
	selector := &fetch.Selector{Allow: allow}
	return fetch.WithSelector(fetch.DefaultFetchClient, selector)
}

func makeValidityURLRule(c *tomlconfig.Config) validity.URLRule {
	return validity.FixedURL(c.SXG.GetValidityURL())
}

func makeProcessor(c *tomlconfig.Config) processor.Processor {
	var tasks []htmltask.HTMLTask

	tasks = append(tasks, htmltask.ConservativeTaskSet...)

	if c.Processor.PreloadCSS {
		tasks = append(tasks, htmltask.PreloadStylesheets())
	}
	if c.Processor.PreloadJS {
		tasks = append(tasks, htmltask.InsecurePreloadScripts())
	}

	config := complexproc.Config{
		Preverify: preverify.Config{MaxContentLength: c.Processor.SizeLimit},
		HTML:      htmlproc.Config{TaskSet: tasks},
	}

	return complexproc.NewComprehensiveProcessor(config)
}

func makeValidPeriodRule(c *tomlconfig.Config) vprule.Rule {
	jsExpiry := c.SXG.GetJSExpiry()
	expiry := c.SXG.GetExpiry()

	return vprule.PerContentType(
		map[string]vprule.Rule{
			"application/javascript":   vprule.FixedLifetime(jsExpiry),
			"application/x-javascript": vprule.FixedLifetime(jsExpiry),
			"text/javascript":          vprule.FixedLifetime(jsExpiry),
		},
		vprule.FixedLifetime(expiry),
	)
}

func makeExchangeFactory(c *tomlconfig.Config) (*ExchangeMetaFactory, error) {
	ec := ExchangeConfig{CertURLBase: c.SXG.GetCertURLBase()}

	var errs *multierror.Error
	var err error

	ec.CertManager, err = makeCertManager(c)
	errs = multierror.Append(errs, err)
	ec.PrivateKey, err = certchainutil.ReadPrivateKeyFile(c.SXG.Cert.KeyFile)
	errs = multierror.Append(errs, err)
	ec.KeepNonSXGPreloads = c.SXG.KeepNonSXGPreloads

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}
	return NewExchangeMetaFactory(ec), nil
}

func makeCertManager(c *tomlconfig.Config) (*certmanager.Manager, error) {
	var rcs certmanager.RawChainSource

	if c.SXG.ACME.Enable {
		key, err := certchainutil.ReadPrivateKeyFile(c.SXG.Cert.KeyFile)
		if err != nil {
			return nil, err
		}

		csr, err := certchainutil.ReadCertificateRequestFile(c.SXG.ACME.CSRFile)
		if err != nil {
			return nil, err
		}

		rcs, err = acmeclient.NewClient(acmeclient.Config{
			CertSignRequest:   csr,
			User:              acmeclient.NewUser(c.SXG.ACME.Email, key),
			DiscoveryURL:      c.SXG.ACME.DiscoveryURL,
			EABHmac:           c.SXG.ACME.EABHmac,
			EABKid:            c.SXG.ACME.EABKid,
			HTTPChallengePort: c.SXG.ACME.HTTPChallengePort,
			HTTPWebRootDir:    c.SXG.ACME.HTTPWebRootDir,
			TLSChallengePort:  c.SXG.ACME.TLSChallengePort,
			DNSProvider:       c.SXG.ACME.DNSProvider,
			ShouldRegister:    true,
			FetchTiming:       certmanager.FetchHourly,
		})
		if err != nil {
			return nil, err
		}
	} else {
		rcs = certmanager.NewLocalCertFile(certmanager.LocalCertFileConfig{
			Path:          c.SXG.Cert.PEMFile,
			AllowTestCert: c.SXG.Cert.AllowTestCert,
		})
	}

	mc := certmanager.Config{
		RawChainSource: rcs,
		OCSPRespSource: certmanager.NewOCSPClient(
			certmanager.OCSPClientConfig{
				AllowTestCert: c.SXG.Cert.AllowTestCert,
			},
		),
	}
	if c.SXG.Cert.CacheDir != "" {
		fmt.Printf("Creating SXG certificate cache directory: %s\n", c.SXG.Cert.CacheDir)
		if err := os.MkdirAll(c.SXG.Cert.CacheDir, 0700); err != nil {
			return nil, err
		}
		// Clean the path with symlinks respected.
		dir, err := filepath.EvalSymlinks(c.SXG.Cert.CacheDir)
		if err != nil {
			return nil, err
		}
		mc.Cache, err = certmanager.NewMultiCertDiskCache(certmanager.MultiCertDiskCacheConfig{
			CertDir:        dir,
			LatestCertFile: "latest.pem",
			LatestOCSPFile: "latest.ocsp",
			LockFile:       ".lock",
		})
		if err != nil {
			return nil, err
		}
	}
	return certmanager.NewManager(mc), nil
}
