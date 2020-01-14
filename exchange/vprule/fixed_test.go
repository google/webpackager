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

package vprule_test

import (
	"testing"
	"time"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/exchange/vprule"
)

func TestFixedLifeTime(t *testing.T) {
	rule := vprule.FixedLifetime(time.Hour)

	got := rule.Get(
		exchangetest.MakeEmptyResponse("https://example.com/dummy/"),
		time.Date(2020, time.January, 10, 20, 15, 0, 0, time.UTC))
	want := exchange.NewValidPeriod(
		time.Date(2020, time.January, 10, 20, 15, 0, 0, time.UTC),
		time.Date(2020, time.January, 10, 21, 15, 0, 0, time.UTC))

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
