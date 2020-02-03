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

// Package preverify implements processors to verify that HTTP responses
// can be distributed as signed exchanges. These processors do not mutate
// the provided exchange.Response: they just inspect it and report an error
// when it does not meet the criteria.
package preverify

import (
	"github.com/google/webpackager/processor"
)

// Config holds the parameters to CheckPrerequisites.
type Config struct {
	// GoodStatusCodes specifies the set of HTTP response codes to consider
	// to be eligible for signed exchanges.
	//
	// nil or empty implies []int{http.StatusOK}, which is considered to be
	// the current best practice.
	GoodStatusCodes []int

	// MaxContentLength specifies the maximum size of each resource turned
	// into a signed exchange, in bytes.
	//
	// Zero implies DefaultMaxContentLength; a negative implies "unlimited."
	MaxContentLength int
}

// The default value(s) used by Config.
const (
	DefaultMaxContentLength = 4194304 // 4 MiB
)

// CheckPrerequisites returns a Processor to verify the provided response
// meets all prerequisites as specified in config.
//
// CheckPrerequisites is usually used indirectly, through the complexproc
// package.
func CheckPrerequisites(config Config) processor.Processor {
	var p processor.SequentialProcessor

	if len(config.GoodStatusCodes) == 0 {
		p = append(p, HTTPStatusOK)
	} else {
		p = append(p, HTTPStatusCode(config.GoodStatusCodes...))
	}

	if config.MaxContentLength >= 0 {
		if config.MaxContentLength == 0 {
			p = append(p, MaxContentLength(DefaultMaxContentLength))
		} else {
			p = append(p, MaxContentLength(config.MaxContentLength))
		}
	}

	return p
}
