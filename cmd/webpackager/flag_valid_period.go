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

package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/google/webpackager/exchange"
	multierror "github.com/hashicorp/go-multierror"
)

var (
	flagDate   = flag.String("date", dateNowString, `Timestamp of signed exchanges in RFC 3339 format ("2006-01-02T15:04:05Z") or "now".`)
	flagExpiry = flag.String("expiry", "1h", `Lifetime of signed exchanges. Maximum is "168h".`)
)

const (
	dateNowString = "now"
	maxExpiry     = 7 * (24 * time.Hour)
)

func getValidPeriodFromFlags() (exchange.ValidPeriod, error) {
	errs := new(multierror.Error)

	date, err := parseDate(*flagDate)
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("invalid --date: %v", err))
	}

	lifetime, err := parseExpiry(*flagExpiry)
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("invalid --expiry: %v", err))
	}

	if err := errs.ErrorOrNil(); err != nil {
		return exchange.ValidPeriod{}, err
	}

	vp := exchange.NewValidPeriodWithLifetime(date, lifetime)
	return vp, nil
}

func parseDate(s string) (time.Time, error) {
	now := time.Now()

	if s == dateNowString {
		return now, nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return now, err
	}
	if t.After(now) {
		return now, errors.New("signing for a future date is disallowed")
	}
	return t, nil
}

func parseExpiry(s string) (time.Duration, error) {
	v, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	if v <= 0 {
		return 0, errors.New("duration must be positive")
	}
	if v > maxExpiry {
		return 0, errors.New("duration too large")
	}
	return v, nil
}
