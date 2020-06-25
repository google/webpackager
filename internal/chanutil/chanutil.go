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

package chanutil

// Killer sends a "kill" signal to goroutines.
type Killer struct {
	// C is a channel to receive the kill signal. It is safe with multiple
	// receive operations from multiple goroutines.
	//
	// Technically, C is just closed when Kill is called. Note that reading
	// from closed channels never gets blocked.
	C <-chan struct{}

	// c is identical to C. We need it for channel closure while we expose
	// a read-only channel for safety.
	c chan struct{}
}

// NewKiller returns a new Killer.
func NewKiller() *Killer {
	c := make(chan struct{})
	return &Killer{
		C: c,
		c: c,
	}
}

// Kill notifies k of the kill signal. Technically, Kill closes k.C to unblock
// all receive operations.
//
// Kill detects the past kill signal and avoids closing k.C multiple times, but
// is unsafe for concurrent calls since the detection and closure is not atomic.
func (k Killer) Kill() {
	select {
	case <-k.c:
		// Already sent a kill signal. Do nothing.
	default:
		close(k.c)
	}
}
