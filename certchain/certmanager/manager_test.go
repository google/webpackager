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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/internal/certchaintest"
)

func TestManager(t *testing.T) {
	augm0409 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/certmanager_0401_0409.cbor")
	augm0413 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/certmanager_0401_0413.cbor")
	augm0415 := certchaintest.MustReadAugmentedChainFile("../../testdata/certs/cbor/certmanager_0415_0415.cbor")

	producer := newStubProducer()
	cache := newStubCache()

	m := certmanager.NewManager(certmanager.Config{
		Producer: producer,
		Cache:    cache,
	})

	producer.out <- augm0409

	if err := m.Start(); err != nil {
		t.Fatalf("m.Start() = error(%q), want success", err)
	}

	{
		t.Log("Started with augm0409")

		select {
		case got := <-cache.OnWrite:
			if diff := cmp.Diff(augm0409, got, certComparer); diff != "" {
				t.Errorf("cache mismatch (-want +got):\n%s", diff)
			}
		case <-time.After(defaultTimeout):
			t.Error("cache timeout")
		}

		got := m.GetAugmentedChain()
		if diff := cmp.Diff(augm0409, got, certComparer); diff != "" {
			t.Errorf("m.GetAugmentedChain() mismatch (-want +got):\n%s", diff)
		}
	}

	producer.out <- augm0413
	{
		t.Log("Updated to augm0413")

		select {
		case got := <-cache.OnWrite:
			if diff := cmp.Diff(augm0413, got, certComparer); diff != "" {
				t.Errorf("cache mismatch (-want +got):\n%s", diff)
			}
		case <-time.After(defaultTimeout):
			t.Error("cache timeout")
		}

		got := m.GetAugmentedChain()
		if diff := cmp.Diff(augm0413, got, certComparer); diff != "" {
			t.Errorf("m.GetAugmentedChain() mismatch (-want +got):\n%s", diff)
		}
	}

	producer.out <- augm0415
	{
		t.Log("Updated to augm0415")

		select {
		case got := <-cache.OnWrite:
			if diff := cmp.Diff(augm0415, got, certComparer); diff != "" {
				t.Errorf("cache mismatch (-want +got):\n%s", diff)
			}
		case <-time.After(defaultTimeout):
			t.Error("cache timeout")
		}

		got := m.GetAugmentedChain()
		if diff := cmp.Diff(augm0415, got, certComparer); diff != "" {
			t.Errorf("m.GetAugmentedChain() mismatch (-want +got):\n%s", diff)
		}
	}
}
