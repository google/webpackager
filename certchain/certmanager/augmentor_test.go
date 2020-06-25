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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"github.com/google/webpackager/internal/certchaintest"
)

func TestFactory(t *testing.T) {
	cert0401 := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0401.pem")
	cert0415 := certchaintest.MustReadRawChainFile("../../testdata/certs/chain/certmanager_0415.pem")
	ocsp0409 := certchaintest.MustReadOCSPRespFile("../../testdata/ocsp/certmanager_0401_0409.ocsp")
	ocsp0413 := certchaintest.MustReadOCSPRespFile("../../testdata/ocsp/certmanager_0401_0413.ocsp")
	ocsp0415 := certchaintest.MustReadOCSPRespFile("../../testdata/ocsp/certmanager_0415_0415.ocsp")
	augm0409 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/certmanager_0401_0409.cbor")
	augm0413 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/certmanager_0401_0413.cbor")
	augm0415 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/certmanager_0415_0415.cbor")

	certSource := newStubRawChainSource(cert0401, cert0401, cert0415)
	ocspSource := newStubOCSPRespSource(ocsp0409, ocsp0413, ocsp0415)

	a := certmanager.NewAugmentor(certSource, ocspSource)

	if err := a.Start(); err != nil {
		t.Fatalf("a.Start() = error(%q), want success", err)
	}

	var certNextRun, ocspNextRun *futureevent.TriggerableEvent

	{
		t.Log("Started with cert0401+ocsp0409")

		if r := waitFor(certSource.OnFetchDone, instantTimeout); r != waitSuccess {
			t.Fatal("certSource wasn't invoked")
		}
		if r := waitFor(ocspSource.OnFetchDone, instantTimeout); r != waitSuccess {
			t.Fatal("ocspSource wasn't invoked")
		}

		var ac *certchain.AugmentedChain
		select {
		case ac = <-a.Out():
			// OK
		default:
			t.Fatal("<-a.Out() empty")
		}
		if diff := cmp.Diff(augm0409, ac, certComparer); diff != "" {
			t.Errorf("<-a.Out() mismatch (-want +got):\n%s", diff)
		}

		if got := certSource.RestCount(); got != 2 {
			t.Fatalf("certSource.RestCount() = %v, want %v", got, 2)
		}
		if got := ocspSource.RestCount(); got != 2 {
			t.Fatalf("ocspSource.RestCount() = %v, want %v", got, 2)
		}

		certNextRun = certSource.NextRun
		ocspNextRun = ocspSource.NextRun
	}

	ocspNextRun.Trigger()
	{
		t.Log("Fetching ocsp0413 (updated)")

		if r := waitFor(ocspSource.OnFetchDone, defaultTimeout); r != waitSuccess {
			t.Fatal("ocspSource wasn't invoked")
		}

		var ac *certchain.AugmentedChain
		select {
		case ac = <-a.Out():
			// OK
		default:
			t.Fatal("<-a.Out() empty")
		}
		if diff := cmp.Diff(augm0413, ac, certComparer); diff != "" {
			t.Errorf("<-a.Out() mismatch (-want +got):\n%s", diff)
		}

		if got := certSource.RestCount(); got != 2 {
			t.Fatalf("certSource.RestCount() = %v, want %v", got, 2)
		}
		if got := ocspSource.RestCount(); got != 1 {
			t.Fatalf("ocspSource.RestCount() = %v, want %v", got, 1)
		}

		certNextRun = certSource.NextRun
		ocspNextRun = ocspSource.NextRun
	}

	certNextRun.Trigger()
	{
		t.Log("Fetching cert0401 (not updated)")

		if r := waitFor(certSource.OnFetchDone, defaultTimeout); r != waitSuccess {
			t.Fatal("certSource wasn't invoked")
		}
		if r := waitFor(ocspSource.OnFetchDone, instantTimeout); r == waitSuccess {
			t.Fatal("ocspSource was invoked unexpectedly")
		}

		select {
		case got := <-a.Out():
			t.Fatalf("<-a.Out() = %#v, want empty", got)
		default:
			// OK
		}

		if got := certSource.RestCount(); got != 1 {
			t.Fatalf("certSource.RestCount() = %v, want %v", got, 1)
		}
		if got := ocspSource.RestCount(); got != 1 {
			t.Fatalf("ocspSource.RestCount() = %v, want %v", got, 1)
		}

		certNextRun = certSource.NextRun
		ocspNextRun = ocspSource.NextRun
	}

	certNextRun.Trigger()
	{
		t.Log("Fetching cert0415 (updated)")

		if r := waitFor(certSource.OnFetchDone, defaultTimeout); r != waitSuccess {
			t.Fatal("certSource wasn't invoked")
		}
		if r := waitFor(ocspSource.OnFetchDone, defaultTimeout); r != waitSuccess {
			t.Fatal("ocspSource wasn't invoked")
		}

		var ac *certchain.AugmentedChain
		select {
		case ac = <-a.Out():
			// OK
		default:
			t.Fatal("<-a.Out() empty")
		}
		if diff := cmp.Diff(augm0415, ac, certComparer); diff != "" {
			t.Errorf("<-a.Out() mismatch (-want +got):\n%s", diff)
		}

		if got := certSource.RestCount(); got != 0 {
			t.Fatalf("certSource.RestCount() = %v, want %v", got, 0)
		}
		if got := ocspSource.RestCount(); got != 0 {
			t.Fatalf("ocspSource.RestCount() = %v, want %v", got, 0)
		}

		if r := waitForEvent(ocspNextRun, defaultTimeout); r != waitCanceled {
			t.Errorf("OCSPSource event hasn't been canceled: %v", r)
		}

		certNextRun = certSource.NextRun
		ocspNextRun = ocspSource.NextRun
	}

	a.Stop()
	{
		t.Log("Stopped")

		if r := waitForEvent(certNextRun, defaultTimeout); r != waitCanceled {
			t.Errorf("CertSource event hasn't been canceled: %v", r)
		}
		if r := waitForEvent(ocspNextRun, defaultTimeout); r != waitCanceled {
			t.Errorf("OCSPSource event hasn't been canceled: %v", r)
		}
	}
}
