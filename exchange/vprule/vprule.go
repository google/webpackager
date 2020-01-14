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

// Package vprule defines how to determine the validity period of signed
// exchanges.
package vprule

import (
	"time"

	"github.com/google/webpackager/exchange"
)

// DefaultRule is the default rule used by webpackager.Packager.
var DefaultRule Rule = FixedLifetime(24 * time.Hour)

// Rule determines the validity period of the provided signed exchange.
type Rule interface {
	// Get returns ValidPeriod for resp. date is the value of vp.Date(),
	// where vp is the ValidPeriod.
	Get(resp *exchange.Response, date time.Time) exchange.ValidPeriod
}
