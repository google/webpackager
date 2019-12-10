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

package exchange

import "time"

// ValidPeriod represents the period the signed exchange is valid for.
type ValidPeriod struct {
	date    time.Time
	expires time.Time
}

// NewValidPeriod creates and initializes a new ValidPeriod from the date and
// expires parameters.
func NewValidPeriod(date, expires time.Time) ValidPeriod {
	return ValidPeriod{date, expires}
}

// NewValidPeriodWithDuration creates and initializes a new ValidPeriod from
// the date parameter and the lifetime.
func NewValidPeriodWithLifetime(date time.Time, lifetime time.Duration) ValidPeriod {
	return ValidPeriod{date, date.Add(lifetime)}
}

// Date returns the date parameter, when the signed exchange is produced.
func (vp ValidPeriod) Date() time.Time {
	return vp.date
}

// Expires returns the expires parameter, when the signed exchange gets expired.
func (vp ValidPeriod) Expires() time.Time {
	return vp.expires
}

// Lifetime returns the duration between the date and expires parameters.
func (vp ValidPeriod) Lifetime() time.Duration {
	return vp.expires.Sub(vp.date)
}

// Contains reports whether t is neither before the date parameter nor after
// the expires parameter. In other words, Contains returns true if t is between
// the date and expires parameters, both inclusive.
func (vp ValidPeriod) Contains(t time.Time) bool {
	return !t.Before(vp.date) && !t.After(vp.expires)
}
