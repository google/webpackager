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
	"os"
	"os/signal"
	"time"
)

// OnSignal returns a new Event kicked on sig.
func OnSignal(sig os.Signal) Event {
	c := make(chan time.Time, 1)
	relay := make(chan os.Signal, 1)

	go func() {
		select {
		case _, ok := <-relay:
			if ok {
				c <- time.Now()
			}
		}
	}()

	signal.Notify(relay, sig)

	return &onSignal{c, relay}
}

type onSignal struct {
	c     chan time.Time
	relay chan os.Signal
}

func (e *onSignal) Chan() <-chan time.Time {
	return e.c
}

func (e *onSignal) Cancel() {
	signal.Stop(e.relay)
	drainOneSignal(e.relay)
	close(e.relay)
	drainOneTime(e.c)
	close(e.c)
}
