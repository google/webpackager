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

package httpheader

import (
	"errors"
	"net/mail"
	"time"
)

var (
	obsoleteLayouts = []string{
		time.RFC850,
		time.ANSIC,
	}

	errUnknownFormat = errors.New("httpdate: unknown format")
)

// ParseDate parses the provided Date header value. ParseDate accepts
// RFC 5322, RFC 850, and ANSI-C asctime formats, which include all formats
// required by RFC 7231.
func ParseDate(value string) (time.Time, error) {
	if t, err := mail.ParseDate(value); err == nil {
		return t, nil
	}

	for _, layout := range obsoleteLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errUnknownFormat
}
