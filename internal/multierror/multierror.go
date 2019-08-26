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

// Package multierror provides interface to accumulate multiple errors and
// handle them as a single error.
package multierror

import (
	"strings"
)

// MultiError represents a collection of errors.
type MultiError struct {
	Errors []error
}

// Add adds err to the MultiError.
func (me *MultiError) Add(err error) {
	me.Errors = append(me.Errors, err)
}

// Err returns nil if the MultiError contains no errors, and returns the
// MultiError itself otherwise.
func (me *MultiError) Err() error {
	if len(me.Errors) > 0 {
		return me
	}
	return nil
}

// Error implements the error built-in interface.
func (me *MultiError) Error() string {
	var sb strings.Builder
	for i, err := range me.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}
