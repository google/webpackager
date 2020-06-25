// Copyright 2020 Google LLC
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

package tomlconfig

import (
	"errors"
	"fmt"
)

// Error reports an error with TOML config.
type Error struct {
	Name string
	Err  error
}

func wrapError(name string, err error) error {
	return &Error{name, err}
}

func newError(name, msg string) *Error {
	return &Error{name, errors.New(msg)}
}

// Error implements error.
func (e *Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Name, e.Err)
}

// Unwrap returns e.Err.
func (e *Error) Unwrap() error { return e.Err }
