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

package exchange_test

import (
	"testing"
	"time"

	"github.com/google/webpackager/exchange"
)

func TestNewValidPeriod(t *testing.T) {
	date := time.Date(2019, time.October, 1, 10, 30, 0, 0, time.UTC)
	expires := time.Date(2019, time.October, 1, 19, 30, 0, 0, time.UTC)

	vp := exchange.NewValidPeriod(date, expires)

	if got := vp.Date(); !got.Equal(date) {
		t.Errorf("vp.Date() = %v, wants %v", got, date)
	}
	if got := vp.Expires(); !got.Equal(expires) {
		t.Errorf("vp.Expires() = %v, wants %v", got, expires)
	}
	if got := vp.Lifetime(); got != 9*time.Hour {
		t.Errorf("vp.Lifetime() = %v, wants %v", got, 9*time.Hour)
	}
}

func TestNewValidPeriodWithLifetime(t *testing.T) {
	date := time.Date(2019, time.October, 1, 10, 30, 0, 0, time.UTC)
	expires := time.Date(2019, time.October, 1, 19, 30, 0, 0, time.UTC)

	vp := exchange.NewValidPeriodWithLifetime(date, 9*time.Hour)

	if got := vp.Date(); !got.Equal(date) {
		t.Errorf("vp.Date() = %v, wants %v", got, date)
	}
	if got := vp.Expires(); !got.Equal(expires) {
		t.Errorf("vp.Expires() = %v, wants %v", got, expires)
	}
	if got := vp.Lifetime(); got != 9*time.Hour {
		t.Errorf("vp.Lifetime() = %v, wants %v", got, 9*time.Hour)
	}
}
