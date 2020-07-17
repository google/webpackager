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

package certmanager

import (
	"bytes"
	"errors"
	"path"

	"github.com/gofrs/flock"
	"github.com/google/renameio"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certchainutil"
	"github.com/hashicorp/go-multierror"
)

var (
	errEmptyCertDir  = errors.New("certmanager: empty CertDir")
	errEmptyCertFile = errors.New("certmanager: empty CertFile")
	errEmptyDigest   = errors.New("certmanager: empty Digest")
	errEmptyLockFile = errors.New("certmanager: empty LockFile")
	errEmptyOCSPFile = errors.New("certmanager: empty OCSPFile")
)

// MultiCertDiskCache is a Cache on a local filesystem. It writes the
// certificate chain in the PEM format and the OCSP response in the DER format
// to separate files as specified by MultiCertDiskCacheConfig. It uses the
// digest of the certificate chain as the basename for the certificate and OCSP
// files.
type MultiCertDiskCache struct {
	MultiCertDiskCacheConfig
}

var _ Cache = (*MultiCertDiskCache)(nil)

// MultiCertDiskCacheConfig configures DiskCache.
type MultiCertDiskCacheConfig struct {
	// CertDir locates the directory to write the certificate chain to.
	// If CertDir is empty, NewMultiCertDiskCache returns an error.
	CertDir string

	// LatestCertFile specifies the filename to be used for the latest
	// version of the certificate. The file will be located in CertDir.
	// If LatestCertFile is empty, NewMultiCertDiskCache returns an error.
	LatestCertFile string

	// LatestOCSPFile specifies the filename to be used for the latest
	// version of the OCSP. The file will be located in CertDir.
	// If LatestOCSPFile is empty, NewMultiCertDiskCache returns an error.
	LatestOCSPFile string

	// LockFile locates the lock file. Must be non-empty. The file will be
	// located in CertDir. If LockFile is empty, NewMultiCertDiskCache returns
	// an error.
	LockFile string
}

// NewMultiCertDiskCache creates and initializes a new MultiCertDiskCache.
func NewMultiCertDiskCache(config MultiCertDiskCacheConfig) (*MultiCertDiskCache, error) {
	var errs *multierror.Error

	if config.CertDir == "" {
		errs = multierror.Append(errs, errEmptyCertDir)
	}
	if config.LockFile == "" {
		errs = multierror.Append(errs, errEmptyLockFile)
	}
	if config.LatestCertFile == "" {
		errs = multierror.Append(errs, errEmptyCertFile)
	}
	if config.LatestOCSPFile == "" {
		errs = multierror.Append(errs, errEmptyOCSPFile)
	}
	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	return &MultiCertDiskCache{config}, nil
}

// Read reads the certificate chain and the OCSP response from local files
// and reproduces an AugmentedChain. Read uses digest as the base filename
// used to retrieve the elements needed to construct the AugmentedChain.
// Read returns a multierror.Error (hashicorp/go-multierror) to report as many
// problems as possible.
func (d *MultiCertDiskCache) Read(digest string) (*certchain.AugmentedChain, error) {
	var errs *multierror.Error

	if digest == "" {
		errs = multierror.Append(errs, errEmptyDigest)
	}

	lock := flock.New(path.Join(d.CertDir, d.LockFile))
	errs = multierror.Append(errs, lock.RLock())

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	cert := path.Join(d.CertDir, digest+".pem")
	ocsp := path.Join(d.CertDir, digest+".ocsp")
	rawChain, err := certchainutil.ReadRawChainFile(cert)
	errs = multierror.Append(errs, err)
	ocspResp, err := certchainutil.ReadOCSPRespFile(ocsp)
	errs = multierror.Append(errs, err)

	var augChain *certchain.AugmentedChain
	if errs.ErrorOrNil() == nil {
		augChain = certchain.NewAugmentedChain(rawChain, ocspResp, nil)
	}

	errs = multierror.Append(errs, lock.Unlock())
	return augChain, errs.ErrorOrNil()
}

// ReadLatest returns the latest version of the AugmentedChain in the cache,
// ErrNotFound otherwise.
func (d *MultiCertDiskCache) ReadLatest() (*certchain.AugmentedChain, error) {
	var errs *multierror.Error

	lock := flock.New(path.Join(d.CertDir, d.LockFile))
	errs = multierror.Append(errs, lock.RLock())

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	cert := path.Join(d.CertDir, d.LatestCertFile)
	ocsp := path.Join(d.CertDir, d.LatestOCSPFile)
	rawChain, err := certchainutil.ReadRawChainFile(cert)
	errs = multierror.Append(errs, err)
	ocspResp, err := certchainutil.ReadOCSPRespFile(ocsp)
	errs = multierror.Append(errs, err)

	var augChain *certchain.AugmentedChain
	if errs.ErrorOrNil() == nil {
		augChain = certchain.NewAugmentedChain(rawChain, ocspResp, nil)
	}

	errs = multierror.Append(errs, lock.Unlock())
	return augChain, errs.ErrorOrNil()
}

// Write writes the certificate chain and the OCSP response from ac into local
// files. It returns a multierror.Error (hashicorp/go-multierror) to report as
// many problems as possible.
func (d *MultiCertDiskCache) Write(ac *certchain.AugmentedChain) error {
	var errs *multierror.Error

	lock := flock.New(path.Join(d.CertDir, d.LockFile))
	errs = multierror.Append(errs, lock.Lock())

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	errs = multierror.Append(errs, d.writeCert(ac))
	errs = multierror.Append(errs, d.writeOCSP(ac))
	errs = multierror.Append(errs, lock.Unlock())

	return errs.ErrorOrNil()
}

func (d *MultiCertDiskCache) writeCert(ac *certchain.AugmentedChain) error {
	buf := new(bytes.Buffer)
	if err := ac.WritePEM(buf); err != nil {
		return err
	}

	certPath := path.Join(d.CertDir, ac.Digest+".pem")
	if err := renameio.WriteFile(certPath, buf.Bytes(), 0600); err != nil {
		return err
	}

	latest := path.Join(d.CertDir, d.LatestCertFile)
	return renameio.WriteFile(latest, buf.Bytes(), 0600)
}

func (d *MultiCertDiskCache) writeOCSP(ac *certchain.AugmentedChain) error {
	ocspPath := path.Join(d.CertDir, ac.Digest+".ocsp")
	if err := renameio.WriteFile(ocspPath, ac.OCSPResp.Raw, 0600); err != nil {
		return err
	}

	latest := path.Join(d.CertDir, d.LatestOCSPFile)
	return renameio.WriteFile(latest, ac.OCSPResp.Raw, 0600)
}
