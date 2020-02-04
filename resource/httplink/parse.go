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

package httplink

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
)

const (
	// https://tools.ietf.org/html/rfc7230#section-3.2.6
	token  = "[!#$%&'*+\\-.^_`|~0-9A-Za-z]+"
	quoted = `"(?:[^"\\]|\\.)*"`
	// https://tools.ietf.org/html/rfc8288#section-3
	uri       = `[^<>]*` // Parsed by url.Parse.
	linkParam = token + `\s*(?:=\s*(?:` + token + `|` + quoted + `))?`
	linkValue = `<(` + uri + `)>(?:\s*;\s*` + linkParam + `)*`
)

var (
	// reValue matches a single "link-value" [RFC 8288].
	reValue = regexp.MustCompile(`^\s*(` + linkValue + `)\s*(?:,|$)`)
	// reParam matches a single "link-param" [RFC 8288].
	reParam = regexp.MustCompile(`;\s*(` + linkParam + `)`)
)

// Parse parses the Link HTTP header value. It returns a slice since the
// value may contain multiple links, separated by comma.
func Parse(header string) ([]*Link, error) {
	var links []*Link
	s := header
	for len(s) > 0 {
		m := reValue.FindStringSubmatch(s)
		if m == nil {
			return nil, fmt.Errorf("malformed Link header: %q", header)
		}
		s = s[len(m[0]):]

		rawLink := m[1] // Captures linkValue.
		rawURL := m[2]  // Captures uri.

		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, xerrors.Errorf("invalid Link URL: %w", err)
		}

		params := make(LinkParams)

		for _, m := range reParam.FindAllStringSubmatch(rawLink, -1) {
			kv := strings.SplitN(m[1], "=", 2)

			key := strings.ToLower(strings.TrimSpace(kv[0]))
			if len(kv) < 2 {
				params.Set(key, "")
			} else {
				val := strings.TrimSpace(kv[1])
				if strings.HasPrefix(val, "\"") {
					val = strings.Replace(val[1:len(val)-1], "\\", "", -1)
				}
				params.Set(key, val)
			}
		}

		links = append(links, &Link{u, params})
	}

	return links, nil
}
