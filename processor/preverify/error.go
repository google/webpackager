// Copyright 2020 Google LLC
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

package preverify

import (
	"fmt"
)

// HTTPStatusError represents an HTTP status error.
type HTTPStatusError struct {
	// StatusCode represents the HTTP status code returned.
	StatusCode int
}

// NewHTTPStatusError creates a new HTTPStatusError.
func NewHTTPStatusError(statusCode int) *HTTPStatusError {
	return &HTTPStatusError{statusCode}
}

// Error implements the error interface.
func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("server responded with status code %d", e.StatusCode)
}
