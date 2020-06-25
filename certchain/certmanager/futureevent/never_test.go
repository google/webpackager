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

func TestNeverOccurs(t *testing.T) {
	e := futureevent.NeverOccurs()
	time.Sleep(time.Microsecond)
	select {
	case _, ok := <-e.Chan():
		if ok {
			t.Error("got notified of the event which should never occur")
		} else {
			t.Error("the notification channel has been closed")
		}
	default:
	}
}

func TestNeverOccursCancel(t *testing.T) {
	e := futureevent.NeverOccurs()
	e.Cancel()
	select {
	case _, ok := <-e.Chan():
		if ok {
			t.Error("the event hasn't been canceled: got notified")
		}
	default:
		t.Error("the event hasn't been canceled: still waiting")
	}
}
