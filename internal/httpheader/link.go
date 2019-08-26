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
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const (
	// https://tools.ietf.org/html/rfc7230#section-3.2.6
	token  = "[!#$%&'*+\\-.^_`|~0-9A-Za-z]+"
	quoted = `"(?:[^"\\]|\\.)*"`

	// https://tools.ietf.org/html/rfc8288#section-3
	linkParam = token + `\s*(?:=\s*(?:` + token + `|` + quoted + `))?`
	linkValue = `<([^<>]*)>(?:\s*;\s*` + linkParam + `)*`
)

var (
	reValue = regexp.MustCompile(`^\s*(` + linkValue + `)\s*(?:,|$)`)
	reParam = regexp.MustCompile(`;\s*(` + linkParam + `)`)
)

// Link represents a parsed Link HTTP header.
type Link struct {
	// URL represents the link target. It can be relative.
	URL *url.URL

	// Params represents the parameters. The keys are lowercased.
	Params map[string]string

	// Header represents the raw header string. It does not contain comma
	// separators. Note it can still contain commas in quoted strings.
	Header string
}

// ParseLink parses a Link header value.
func ParseLink(value string) ([]*Link, error) {
	links := []*Link{}

	for rest := value; len(rest) > 0; {
		m := reValue.FindStringSubmatch(rest)
		if m == nil {
			return nil, fmt.Errorf("malforemd Link header: %s", value)
		}
		rest = rest[len(m[0]):]

		linkValue := m[1]
		rawurl := m[2]

		u, err := url.Parse(rawurl)
		if err != nil {
			return nil, fmt.Errorf("malformed Link URL: %v", err)
		}

		params := make(map[string]string)

		for _, m := range reParam.FindAllStringSubmatch(linkValue, -1) {
			kv := strings.SplitN(m[1], "=", 2)

			key := strings.ToLower(strings.TrimSpace(kv[0]))
			if len(kv) < 2 {
				params[key] = ""
			} else {
				val := strings.TrimSpace(kv[1])
				if strings.HasPrefix(val, "\"") {
					val = strings.ReplaceAll(val[1:len(val)-1], "\\", "")
				}
				params[key] = val
			}
		}
		links = append(links, &Link{u, params, linkValue})
	}

	return links, nil
}

// GoString implements the GoStringer interface.
func (link *Link) GoString() string {
	return fmt.Sprintf("&httpheader.Link{URL:&%#v, Params:%#v, Header:%#q}",
		*link.URL, link.Params, link.Header)
}
