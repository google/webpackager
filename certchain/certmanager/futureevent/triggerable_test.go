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

	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

func TestTriggerableEvent(t *testing.T) {
	e := futureevent.NewTriggerableEvent()
	e.Trigger()
	select {
	case _, ok := <-e.Chan():
		if !ok {
			t.Error("the notification channel has been closed")
		}
	default:
		t.Error("got no notification")
	}
}

func TestTriggerableEventNotBeforeTrigger(t *testing.T) {
	e := futureevent.NewTriggerableEvent()
	select {
	case _, ok := <-e.Chan():
		if ok {
			t.Error("got notified of the event before triggering")
		} else {
			t.Error("the notification channel has been closed")
		}
	default:
	}
}

func TestTriggerableEventCancel(t *testing.T) {
	e := futureevent.NewTriggerableEvent()
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
