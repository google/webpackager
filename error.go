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

package webpackager

import (
	"fmt"
	"net/url"
)

// Error represents an error from Packager.Run.
type Error struct {
	// Err represents the actual error.
	Err error
	// URL represents the URL that caused this Error.
	URL *url.URL
}

// WrapError wraps err into an Error. url is the URL which err was raised for.
func WrapError(err error, url *url.URL) error {
	if err == nil {
		return nil
	}
	return &Error{err, url}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("error with processing %s: %v", e.URL, e.Err)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error { return e.Err }
