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

package futureevent_test

import (
	"testing"
	"time"

	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

// In order to reduce the flakiness of the tests which deal with time,
// event scheduling here follows the rules below:
// - If event A is expected to happen after event B, set A 1000 times
//   later than B.
// - For events expected to happen within a reasonable time, wait for
//   up to defaultTimeout (in consts_test.go).

func TestRealClockEvent(t *testing.T) {
	e := futureevent.NewRealClockEvent(time.Now().Add(time.Microsecond))

	select {
	case _, ok := <-e.Chan():
		if !ok {
			t.Error("the notification channel has been closed")
		}
	case <-time.After(defaultTimeout):
		t.Error("timeout")
	}
}

func TestRealClockEventNotTooEarly(t *testing.T) {
	e := futureevent.NewRealClockEvent(time.Now().Add(3 * time.Second))
	time.Sleep(time.Microsecond)
	select {
	case _, ok := <-e.Chan():
		if ok {
			t.Error("got notified of the event too early")
		} else {
			t.Error("the notification channel has been closed")
		}
	default:
	}
}

func TestRealClockEventCancel(t *testing.T) {
	after := 100 * time.Millisecond
	e := futureevent.NewRealClockEvent(time.Now().Add(after))
	e.Cancel()
	time.Sleep(2 * after)
	select {
	case _, ok := <-e.Chan():
		if ok {
			t.Error("the event hasn't been canceled: got notified")
		}
	default:
		t.Error("the event hasn't been canceled: still waiting")
	}
}
