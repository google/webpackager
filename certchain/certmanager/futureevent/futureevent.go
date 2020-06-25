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

// Package futureevent defines interface to handle future events.
package futureevent

import (
	"time"
)

// Event represents a cancelable event that occurs sometime in future.
type Event interface {
	// Chan returns a channel which this Event notifies of the occurrence.
	// The time.Time value represents the time of the occurrence.
	//
	// Remember that the event may be canceled, in which case the channel
	// gets closed and empty. The listeners can detect the cancelation by
	// checking the boolean result of a receive operation.
	Chan() <-chan time.Time

	// Cancel cancels the event.
	//
	// Cancel clears and closes the event notifier channel, thus must not be
	// called more than once.
	Cancel()
}

// Factory creates a new Event raised when the provided time t comes.
type Factory func(t time.Time) Event

// DefaultFactory uses the real clock.
var DefaultFactory Factory = NewRealClockEvent
