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

package futureevent

import (
	"time"
)

// NewRealClockEvent returns a new Event kicked at the provided Time,
// using the real clock.
func NewRealClockEvent(at time.Time) Event {
	c := make(chan time.Time, 1)
	t := time.AfterFunc(at.Sub(time.Now()), func() { c <- time.Now() })
	return &realClockEvent{c, t}
}

type realClockEvent struct {
	// c works as a clearable and closable proxy for t.C. Event requires
	// Chan to be empty and closed once the event has been canceled, but
	// time.Timer does not close or clear its channel on Stop, nor allow us
	// to close it by ourselves (time.Timer.C is read-only).
	c chan time.Time
	t *time.Timer
}

func (e *realClockEvent) Chan() <-chan time.Time {
	return e.c
}

func (e *realClockEvent) Cancel() {
	e.t.Stop()
	drainOneTime(e.t.C)
	drainOneTime(e.c)
	close(e.c)
}
