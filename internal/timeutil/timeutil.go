// Copyright 2020 Google LLC
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

// Package timeutil provides time.Now that can be monkey-patched.
package timeutil

import "time"

var now func() time.Time = time.Now

// Now returns the current local time. It is equal to time.Now by default
// but can be faked using StubNow or other Stub functions.
func Now() time.Time { return now() }

// ResetNow restores the Now's original behavior.
func ResetNow() { now = time.Now }

// StubNow substitutes Now with stub for testing.
func StubNow(stub func() time.Time) { now = stub }

// StubNowWithFixedTime makes Now return newNow for testing. The stubbed
// time does not advance.
func StubNowWithFixedTime(newNow time.Time) {
	StubNow(func() time.Time { return newNow })
}

// StubNowToAdjust "adjusts" Now to newNow for testing. The stubbed time
// continues to advance: Now will return newNow plus the time elapsed since
// StubNowToAdjust is called.
func StubNowToAdjust(newNow time.Time) {
	delta := newNow.Sub(time.Now())
	StubNow(func() time.Time { return time.Now().Add(delta) })
}
