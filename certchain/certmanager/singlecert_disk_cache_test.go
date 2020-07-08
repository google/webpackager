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
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/internal/certchaintest"
)

func TestDiskCache(t *testing.T) {
	p256 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/ecdsap256_nosct.cbor")
	p384 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/ecdsap384_nosct.cbor")

	wants := []struct {
		certSHA1, ocspSHA1 string
	}{
		{
			// ../../testdata/certs/chain/ecdsap256.pem
			certSHA1: "755ccd028cd2a691bd46aff8a90bc216aafd7be2",
			// ../../testdata/ocsp/ecdsap256_7days.ocsp
			ocspSHA1: "0d5debcc3049fa980b3c3626440a7f51e68a5b85",
		},
		{
			// ../../testdata/certs/chain/ecdsap384.pem
			certSHA1: "aa8d28fce08cd4505219f3ce1c66115af144a316",
			// ../../testdata/ocsp/ecdsap384_7days.ocsp
			ocspSHA1: "c08907bf0bb98be23d663e8a523c36c0317341d6",
		},
	}

	tempDir, err := ioutil.TempDir("", "certmanager_test_")
	if err != nil {
		t.Fatalf("cannot set up a test directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "cert.pem")
	ocspPath := filepath.Join(tempDir, "ocsp.der")
	lockPath := filepath.Join(tempDir, ".lock")

	d := certmanager.NewSingleCertDiskCache(certmanager.SingleCertDiskCacheConfig{
		CertPath: certPath,
		OCSPPath: ocspPath,
		LockPath: lockPath,
	})

	if err := d.Write(p256); err != nil {
		t.Fatalf("failed to initialize the cache: %v", err)
	}

	var wg sync.WaitGroup

	// writer writes to the cache many times.
	writer := func() {
		const numLoop = 10
		defer wg.Done()

		for i := 1; i <= numLoop; i++ {
			if err := d.Write(p256); err != nil {
				t.Errorf("iteration #%v: cannot write p256: %v", i, err)
			}
			if err := d.Write(p384); err != nil {
				t.Errorf("iteration #%v: cannot write p384: %v", i, err)
			}
		}
	}

	// reader reads a cert/ocsp pair from the cache and verifies it.
	reader := func(id int) {
		defer wg.Done()

		ac, err := d.ReadLatest()
		if err != nil {
			t.Errorf("reader #%v: error with d.Read(): %v", id, err)
			return
		}

		certSHA1 := sha1.New()
		ac.WritePEM(certSHA1)
		ocspSHA1 := sha1.New()
		ocspSHA1.Write(ac.OCSPResp.Raw)

		got := struct{ certSHA1, ocspSHA1 string }{
			fmt.Sprintf("%x", certSHA1.Sum(nil)),
			fmt.Sprintf("%x", ocspSHA1.Sum(nil)),
		}

		for _, want := range wants {
			if got == want {
				return
			}
		}

		t.Errorf("reader #%v: %+v, want any of %+v", id, got, wants)
	}

	const numReaders = 20
	wg.Add(1 + numReaders)
	go writer()
	for i := 1; i <= numReaders; i++ {
		go reader(i)
	}
	wg.Wait()
}
