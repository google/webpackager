// Copyright 2019 Google LLC
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

package certmanager

import (
	"os"
	"time"

	"github.com/google/webpackager/certchain/certmanager/futureevent"
	"github.com/google/webpackager/internal/timeutil"
)

// FetchHourly makes Fetch called every hour.
var FetchHourly FetchTiming = FetchAtIntervals(time.Hour)

// FetchTiming controls the frequency of the Fetch calls on RawChainSource
// or OCSPRespSource.
type FetchTiming interface {
	// GetNextRun determines the nextRun return parameter.
	GetNextRun() futureevent.Event
}

// FetchTimingFunc turns a function into a FetchTiming.
type FetchTimingFunc func() futureevent.Event

// GetNextRun calls f().
func (f FetchTimingFunc) GetNextRun() futureevent.Event {
	return f()
}

// FetchAtIntervals makes Fetch called at fixed intervals.
func FetchAtIntervals(interval time.Duration) FetchTiming {
	return &fetchAtIntervals{interval, futureevent.DefaultFactory}
}

// FetchAtIntervalsWithEventFactory is like FetchAtIntervals but uses factory
// instead of futureevent.DefaultFactory.
func FetchAtIntervalsWithEventFactory(interval time.Duration, factory futureevent.Factory) FetchTiming {
	return &fetchAtIntervals{interval, factory}
}

type fetchAtIntervals struct {
	interval   time.Duration
	newEventAt futureevent.Factory
}

func (f *fetchAtIntervals) GetNextRun() futureevent.Event {
	return f.newEventAt(timeutil.Now().Add(f.interval))
}

// FetchOnSignal makes Fetch called when signaled by sig.
func FetchOnSignal(sig os.Signal) FetchTiming {
	return FetchTimingFunc(func() futureevent.Event {
		return futureevent.OnSignal(sig)
	})
}

// FetchOnlyOnce makes Fetch called only once, not repeatedly.
func FetchOnlyOnce() FetchTiming {
	return FetchTimingFunc(func() futureevent.Event {
		return futureevent.NeverOccurs()
	})
}
