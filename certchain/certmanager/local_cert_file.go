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
	"time"

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certchainutil"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

// LocalCertFile is a RawChainSource which reads the certificate chain from
// a local file in the PEM format.
type LocalCertFile struct {
	LocalCertFileConfig
}

type LocalCertFileConfig struct {
	// Path locates the PEM file containing the certificate chain.
	Path string

	// FetchTiming controls the frequency of checking for the certificate.
	// nil implies certmanager.FetchHourly.
	FetchTiming FetchTiming

	// AllowTestCert specifies whether to allow test certificates.
	//
	// LocalCertFile calls VerifyChain and VerifySXGCriteria to make sure
	// RawChain is valid for use with signed exchanges. If AllowTestCert
	// is set true, LocalCertFile skips VerifySXGCriteria and accepts any
	// RawChain as long as it is valid in terms of VerifyChain.
	AllowTestCert bool
}

var _ RawChainSource = (*LocalCertFile)(nil)

// NewLocalCertFile creates and initializes a new LocalCertFile.
func NewLocalCertFile(c LocalCertFileConfig) *LocalCertFile {
	if c.FetchTiming == nil {
		c.FetchTiming = FetchHourly
	}
	return &LocalCertFile{c}
}

// Fetch reads the certificate chain from l.Path and returns it as newChain,
// whether or not it is updated from chain.
func (l *LocalCertFile) Fetch(chain *certchain.RawChain, now func() time.Time) (newchain *certchain.RawChain, nextRun futureevent.Event, err error) {
	c, err := certchainutil.ReadRawChainFile(l.Path)

	if err != nil {
		return nil, l.FetchTiming.GetNextRun(), err
	}

	if err := c.VerifyChain(now()); err != nil {
		return nil, l.FetchTiming.GetNextRun(), err
	}
	if err := c.VerifySXGCriteria(); !l.AllowTestCert && err != nil {
		return nil, l.FetchTiming.GetNextRun(), err
	}

	return c, l.FetchTiming.GetNextRun(), nil
}
