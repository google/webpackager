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

package acmeclient

import (
	"crypto/x509"
	"errors"
	"strconv"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/challenge/http01"
	"github.com/go-acme/lego/v3/challenge/tlsalpn01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/providers/http/webroot"
	"github.com/go-acme/lego/v3/registration"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"golang.org/x/xerrors"
)

// Client acquires a signed exchange certificate using the ACME protocol
// as a certmanager.RawChainSource.
type Client struct {
	LegoClient      *lego.Client
	CertSignRequest *x509.CertificateRequest
	FetchTiming     certmanager.FetchTiming
}

var _ certmanager.RawChainSource = (*Client)(nil)

// Config configures Client. It contains information to request a certificate
// to the Certificate Authority (CA) using the ACME protocol.
//
// HTTPChallengePort, HTTPWebRootDir, TLSChallengePort and DNSProvider specify
// how Client responds to challenges from the ACME server. The ACME standard
// defines three types of challenges, namely HTTP, DNS and TLS, and each field
// configures one of them. Only one of the four fields is expected to be set.
// For wildcard domains, the DNS challenge is the only option thus DNSProvider
// must be set.
//
// https://letsencrypt.org/docs/challenge-types/ describes these challenges
// in greater detail.
//
// Port Usage
//
// Client uses HTTPChallengePort or TLSChallengePort, while the ACME protocol
// requires the HTTP and TLS challenge responders to listen on the standard
// ports (80 and 443), because webpackager isn't designed to run as root thus
// can't bind a listener to the privileged ports. Keep in mind that you need
// to proxy challenge traffic to the custom port you specified.
type Config struct {
	// CertSignRequest is the Certificate Signing Request (CSR) sent over the
	// ACME protocol.
	CertSignRequest *x509.CertificateRequest

	// User provides the user information for the ACME request.
	User *User

	// DiscoveryURL is the Discovery Resource URL provided by the Certificate
	// Authority to make ACME requests.
	DiscoveryURL string

	// HTTPChallengePort is the HTTP challenge port used for the ACME HTTP
	// challenge.
	//
	// Remember you need to proxy challenge traffic. See Port Usage above.
	HTTPChallengePort int

	// HTTPWebRootDir is the web root directory where the ACME HTTP challenge
	// token will be deposited.
	HTTPWebRootDir string

	// TLSChallengePort is the TLS challenge port used for the ACME TLS
	// challenge.
	//
	// Remember you need to proxy challenge traffic. See Port Usage above.
	TLSChallengePort int

	// DNSProvider is the ACME DNS Provider used for the challenge. It is
	// specified by the Lego config code.
	//
	// The binary must be built with `-tags dns01` to use the DNS challenge.
	// If DNSProvider is set non-empty in a binary without that build option,
	// NewClient will fail with an error.
	//
	// See https://go-acme.github.io/lego/dns/ for the DNS Provider list.
	DNSProvider string

	// ShouldRegister specifies whether the ACME user needs to register with
	// the Certificate Authority.
	ShouldRegister bool

	// FetchTiming controls the frequency of checking for the certificate.
	// nil implies certmanager.FetchHourly.
	FetchTiming certmanager.FetchTiming
}

// NewClient creates and initializes a new Client with config.
func NewClient(config Config) (*Client, error) {
	legoConfig := lego.NewConfig(config.User)

	legoConfig.CADirURL = config.DiscoveryURL
	legoConfig.Certificate.KeyType = certcrypto.EC256

	legoClient, err := lego.NewClient(legoConfig)
	if err != nil {
		return nil, xerrors.Errorf("obtaining LEGO client: %w", err)
	}

	if config.HTTPChallengePort != 0 {
		s := http01.NewProviderServer("", strconv.Itoa(config.HTTPChallengePort))
		if err := legoClient.Challenge.SetHTTP01Provider(s); err != nil {
			return nil, xerrors.Errorf("setting up HTTP01 challenge provider: %w", err)
		}
	}
	if config.HTTPWebRootDir != "" {
		httpProvider, err := webroot.NewHTTPProvider(config.HTTPWebRootDir)
		if err != nil {
			return nil, xerrors.Errorf("getting HTTP01 challenge provider: %w", err)
		}
		if err := legoClient.Challenge.SetHTTP01Provider(httpProvider); err != nil {
			return nil, xerrors.Errorf("setting up HTTP01 challenge provider: %w", err)
		}
	}

	if config.TLSChallengePort != 0 {
		s := tlsalpn01.NewProviderServer("", strconv.Itoa(config.TLSChallengePort))
		if err := legoClient.Challenge.SetTLSALPN01Provider(s); err != nil {
			return nil, xerrors.Errorf("setting up TLSALPN01 challenge provider: %w", err)
		}
	}

	if config.DNSProvider != "" {
		provider, err := dnsProvider(config.DNSProvider)
		if err != nil {
			return nil, xerrors.Errorf("getting DNS01 challenge provider: %w", err)
		}
		if err := legoClient.Challenge.SetDNS01Provider(provider); err != nil {
			return nil, xerrors.Errorf("setting up DNS01 challenge provider: %w", err)
		}
	}

	// Theoretically, this should always be set to false as users should have
	// pre-registered for access to the ACME CA and agreed to the Terms of
	// Service.
	// TODO(banaag): revisit this when trying the class out with DigiCert CA.
	if !config.ShouldRegister {
		config.User.SetRegistration(new(registration.Resource))
	} else {
		// TODO(banaag): make sure we present the TOS URL to the user and
		// prompt for confirmation.
		// The plan is to move this to some separate setup command outside the
		// server which would be executed one time. Alternatively, we can have
		// a field in the toml file that is documented to indicate agreement
		// with TOS.
		reg, err := legoClient.Registration.Register(
			registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return nil, xerrors.Errorf("ACME CA client registration: %w", err)
		}
		config.User.SetRegistration(reg)
	}

	fetchTiming := config.FetchTiming
	if fetchTiming == nil {
		fetchTiming = certmanager.FetchHourly
	}

	return &Client{
		LegoClient:      legoClient,
		CertSignRequest: config.CertSignRequest,
		FetchTiming:     fetchTiming,
	}, nil
}

// Fetch acquires a new RawChain from the ACME server.
func (c *Client) Fetch(chain *certchain.RawChain, now func() time.Time) (newChain *certchain.RawChain, nextRun futureevent.Event, err error) {
	// Each resource comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL.
	resource, err := c.LegoClient.Certificate.ObtainForCSR(*c.CertSignRequest, true)
	if err != nil {
		return nil, c.FetchTiming.GetNextRun(), err
	}

	if resource == nil || resource.Certificate == nil {
		err = errors.New("acmeclient: no certificate returned")
		return nil, c.FetchTiming.GetNextRun(), err
	}

	newChain, err = certchain.NewRawChainFromPEM(resource.Certificate)
	if err != nil {
		return nil, c.FetchTiming.GetNextRun(), err
	}

	if err := newChain.VerifyChain(now()); err != nil {
		return nil, c.FetchTiming.GetNextRun(), err
	}
	if err := newChain.VerifySXGCriteria(); err != nil {
		return nil, c.FetchTiming.GetNextRun(), err
	}

	return newChain, c.FetchTiming.GetNextRun(), err
}
