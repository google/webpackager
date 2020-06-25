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

package chanutil_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/webpackager/internal/chanutil"
)

func TestKiller(t *testing.T) {
	k := chanutil.NewKiller()

	select {
	case <-k.C:
		t.Fatal("got kill signal before Kill called")
	default:
		// OK
	}

	k.Kill()
	time.Sleep(time.Microsecond)

	select {
	case <-k.C:
		// OK
	default:
		t.Error("coudn't receive a kill signal after Kill called")
	}
}

func TestKillerSecondKill(t *testing.T) {
	k := chanutil.NewKiller()
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("got %v, want success", err)
		}
	}()

	k.Kill()
	k.Kill()
}

func TestKillerMultipleReceives(t *testing.T) {
	const timeout = 3 * time.Second
	const numGoroutines = 8

	k := chanutil.NewKiller()
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			t.Helper()

			select {
			case <-k.C:
				// OK
			case <-time.After(timeout):
				t.Errorf("timeout in goroutine #%d", i)
			}
		}(i)
	}
	k.Kill()
	wg.Wait()
}
