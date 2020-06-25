// Copyright 2020 Google LLC
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

// NeverOccursEvent never occurs. It is still cancelable.
type NeverOccursEvent struct {
	c chan time.Time
}

// NeverOccurs creates a new NeverOccursEvent.
func NeverOccurs() *NeverOccursEvent {
	return &NeverOccursEvent{make(chan time.Time)}
}

// Chan returns an event notifier channel. Since the event never occurs,
// the receive operation on the returned channel would be blocked until
// e is canceled.
func (e *NeverOccursEvent) Chan() <-chan time.Time { return e.c }

// Cancel unblocks the receive operation on e.Chan().
func (e *NeverOccursEvent) Cancel() { close(e.c) }
