// Copyright 2019 Google LLC
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

package certmanager

import (
	"bytes"
	"errors"

	"github.com/gofrs/flock"
	"github.com/google/renameio"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certchainutil"
	"github.com/hashicorp/go-multierror"
)

var (
	errEmptyCertPath = errors.New("certmanager: empty CertPath")
	errEmptyOCSPPath = errors.New("certmanager: empty OCSPPath")
)

// SingleCertDiskCache is a Cache on a local filesystem. It writes the
// certificate chain in the PEM format and the OCSP response in the DER format
// to separate files as specified by SingleCertDiskCacheConfig.
type SingleCertDiskCache struct {
	SingleCertDiskCacheConfig
}

var _ Cache = (*SingleCertDiskCache)(nil)

// SingleCertDiskCacheConfig configures SingleCertDiskCache.
type SingleCertDiskCacheConfig struct {
	// CertPath locates the PEM file to write the certificate chain to.
	// If CertPath is empty, the certificate chain is not cached.
	CertPath string

	// OCSPPath locates the file to write the OCSP response DER bytes to.
	// If OCSPPath is empty, the OCSP response is not cached.
	OCSPPath string

	// LockPath locates the lock file. Must be non-empty.
	LockPath string
}

// NewSingleCertDiskCache creates and initializes a new SingleCertDiskCache.
func NewSingleCertDiskCache(config SingleCertDiskCacheConfig) *SingleCertDiskCache {
	return &SingleCertDiskCache{config}
}

// ReadLatest returns the cached AugmentedChain, otherwise it returns an error.
func (d *SingleCertDiskCache) ReadLatest() (*certchain.AugmentedChain, error) {
	var errs *multierror.Error

	if d.CertPath == "" {
		errs = multierror.Append(errs, errEmptyCertPath)
	}
	if d.OCSPPath == "" {
		errs = multierror.Append(errs, errEmptyOCSPPath)
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	lock := flock.New(d.LockPath)
	errs = multierror.Append(errs, lock.RLock())

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	rawChain, err := certchainutil.ReadRawChainFile(d.CertPath)
	if err != nil {
		errs = multierror.Append(errs, ErrNotFound)
		errs = multierror.Append(errs, err)
	}
	ocspResp, err := certchainutil.ReadOCSPRespFile(d.OCSPPath)
	if err != nil {
		errs = multierror.Append(errs, ErrNotFound)
		errs = multierror.Append(errs, err)
	}

	var augChain *certchain.AugmentedChain
	if errs.ErrorOrNil() == nil {
		augChain = certchain.NewAugmentedChain(rawChain, ocspResp, nil)
	}

	errs = multierror.Append(errs, lock.Unlock())
	return augChain, errs.ErrorOrNil()
}

// Read reads the certificate chain and the OCSP response from local files
// and reproduces an AugmentedChain. Read works only when d.CertPath and
// d.OCSPPath are both non-empty and otherwise returns an error. Read returns
// a multierror.Error (hashicorp/go-multierror) to report as many problems as
// possible.
func (d *SingleCertDiskCache) Read(digest string) (*certchain.AugmentedChain, error) {
	// For single cert implementation, digest is not used in reading, just for
	// verification.
	latest, err := d.ReadLatest()
	if err != nil {
		return nil, err
	}
	if latest.Digest != digest {
		return nil, ErrNotFound
	}
	return latest, nil
}

// Write writes the certificate chain and the OCSP response from ac into local
// files. It returns a multierror.Error (hashicorp/go-multierror) to report as
// many problems as possible.
func (d *SingleCertDiskCache) Write(ac *certchain.AugmentedChain) error {
	if d.CertPath == "" && d.OCSPPath == "" {
		return nil
	}

	var errs *multierror.Error

	lock := flock.New(d.LockPath)
	errs = multierror.Append(errs, lock.Lock())

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	errs = multierror.Append(errs, d.writeCert(ac))
	errs = multierror.Append(errs, d.writeOCSP(ac))
	errs = multierror.Append(errs, lock.Unlock())

	return errs.ErrorOrNil()
}

func (d *SingleCertDiskCache) writeCert(ac *certchain.AugmentedChain) error {
	if d.CertPath == "" {
		return nil
	}

	buf := new(bytes.Buffer)
	if err := ac.WritePEM(buf); err != nil {
		return err
	}

	return renameio.WriteFile(d.CertPath, buf.Bytes(), 0600)
}

func (d *SingleCertDiskCache) writeOCSP(ac *certchain.AugmentedChain) error {
	if d.OCSPPath == "" {
		return nil
	}

	return renameio.WriteFile(d.OCSPPath, ac.OCSPResp.Raw, 0600)
}
