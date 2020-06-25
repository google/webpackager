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

package certmanager_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/internal/certchaintest"
	"github.com/google/webpackager/internal/timeutil"
)

func TestLocalCertFile(t *testing.T) {
	tests := []struct {
		name     string
		oldChain *certchain.RawChain
		file     string
	}{
		{
			name:     "NoOldChain",
			oldChain: nil,
			file:     "../../testdata/certs/chain/certmanager_0415.pem",
		},
		{
			name:     "Updated",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/certs/chain/certmanager_0415.pem",
		},
		{
			name:     "NotUpdated",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/certs/chain/certmanager_0401.pem",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFile := createTempFile()
			defer os.Remove(tempFile)

			c := certmanager.LocalCertFileConfig{
				Path: tempFile,
				FetchTiming: certmanager.FetchAtIntervalsWithEventFactory(
					time.Hour,
					newDummyFutureEvent,
				),
			}
			l := certmanager.NewLocalCertFile(c)

			copyFile(test.file, tempFile)
			wantNewChain := certchaintest.MustReadRawChainFile(test.file)

			now := time.Date(2020, time.April, 15, 15, 0, 0, 0, time.UTC)
			wantNextRun := newDummyFutureEvent(
				time.Date(2020, time.April, 15, 16, 0, 0, 0, time.UTC),
			)
			timeutil.StubNowWithFixedTime(now)
			defer timeutil.ResetNow()
			newChain, nextRun, err := l.Fetch(test.oldChain, func() time.Time { return now })

			// nextRun must be valid even if (err != nil).
			if diff := cmp.Diff(wantNextRun, nextRun); diff != "" {
				t.Errorf("nextRun mismatch (-want +got):\n%s", diff)
			}
			if err != nil {
				t.Fatalf("l.Fetch() = error(%q), want success", err)
			}
			if diff := cmp.Diff(wantNewChain, newChain); diff != "" {
				t.Errorf("l.Fetch mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLocalCertFileError(t *testing.T) {
	tests := []struct {
		name     string
		oldChain *certchain.RawChain
		file     string
	}{
		{
			name:     "InvalidFile",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/keys/ecdsap256.key",
		},
		{
			name:     "NonSXGCert",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/certs/chain/non_sxg_cert.pem",
		},
		{
			name:     "LastingTooLong",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/certs/chain/lasting_365days.pem",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFile := createTempFile()
			defer os.Remove(tempFile)

			c := certmanager.LocalCertFileConfig{
				Path: tempFile,
				FetchTiming: certmanager.FetchAtIntervalsWithEventFactory(
					time.Hour,
					newDummyFutureEvent,
				),
			}
			l := certmanager.NewLocalCertFile(c)

			copyFile(test.file, tempFile)

			now := time.Date(2020, time.April, 15, 15, 0, 0, 0, time.UTC)
			wantNextRun := newDummyFutureEvent(
				time.Date(2020, time.April, 15, 16, 0, 0, 0, time.UTC),
			)
			timeutil.StubNowWithFixedTime(now)
			defer timeutil.ResetNow()
			newChain, nextRun, err := l.Fetch(test.oldChain, func() time.Time { return now })

			// nextRun must be valid even if (err != nil).
			if diff := cmp.Diff(wantNextRun, nextRun); diff != "" {
				t.Errorf("nextRun mismatch (-want +got):\n%s", diff)
			}
			if err == nil {
				t.Errorf("l.Fetch() = %#v (success), want error", newChain)
			}
		})
	}
}

func TestLocalCertFileAllowTestCert(t *testing.T) {
	tests := []struct {
		name     string
		oldChain *certchain.RawChain
		file     string
	}{
		{
			name:     "NonSXGCert",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/certs/chain/non_sxg_cert.pem",
		},
		{
			name:     "LastingTooLong",
			oldChain: certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem"),
			file:     "../../testdata/certs/chain/lasting_365days.pem",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFile := createTempFile()
			defer os.Remove(tempFile)

			c := certmanager.LocalCertFileConfig{
				AllowTestCert: true,
				Path:          tempFile,
				FetchTiming: certmanager.FetchAtIntervalsWithEventFactory(
					time.Hour,
					newDummyFutureEvent,
				),
			}
			l := certmanager.NewLocalCertFile(c)

			copyFile(test.file, tempFile)
			wantNewChain := certchaintest.MustReadRawChainFile(test.file)

			now := time.Date(2020, time.April, 15, 15, 0, 0, 0, time.UTC)
			wantNextRun := newDummyFutureEvent(
				time.Date(2020, time.April, 15, 16, 0, 0, 0, time.UTC),
			)
			timeutil.StubNowWithFixedTime(now)
			defer timeutil.ResetNow()
			newChain, nextRun, err := l.Fetch(test.oldChain, func() time.Time { return now })

			// nextRun must be valid even if (err != nil).
			if diff := cmp.Diff(wantNextRun, nextRun); diff != "" {
				t.Errorf("nextRun mismatch (-want +got):\n%s", diff)
			}
			if err != nil {
				t.Fatalf("l.Fetch() = error(%q), want success", err)
			}
			if diff := cmp.Diff(wantNewChain, newChain); diff != "" {
				t.Errorf("l.Fetch mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
