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

// TriggerableEvent is an Event triggered by a method call.
type TriggerableEvent struct {
	c chan time.Time
}

// NewTriggerableEvent creates a new TriggerableEvent.
func NewTriggerableEvent() *TriggerableEvent {
	return &TriggerableEvent{make(chan time.Time, 1)}
}

// Trigger triggers the event.
func (e *TriggerableEvent) Trigger() {
	e.c <- time.Now()
}

// Chan returns the event notifier channel.
func (e *TriggerableEvent) Chan() <-chan time.Time {
	return e.c
}

// Cancel makes e no longer triggerable.
func (e *TriggerableEvent) Cancel() {
	drainOneTime(e.c)
	close(e.c)
}
