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
)

var (
	flagDate = flag.String("date", dateNowString, `Timestamp of signed exchanges in RFC 3339 format ("2006-01-02T15:04:05Z") or "now".`)
)

const (
	dateNowString = "now"
)

func getDateFromFlags() (time.Time, error) {
	date, err := parseDate(*flagDate)
	if err != nil {
		return date, fmt.Errorf("invalid --date: %v", err)
	}
	return date, nil
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
