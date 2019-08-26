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

package validity

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/webpackager/internal/httpheader"
	"github.com/google/webpackager/internal/urlutil"
)

// ValidityURLRule decides the validity URL of a resource.
type ValidityURLRule interface {
	// Apply returns the validity URL of a resource. physurl and resp are
	// the physical URL and the HTTP response of the resource respectively.
	// The physical URL is typically equal to the request URL, but different
	// in some cases; see package urlrewrite for more details.
	//
	// Note ValidityURLRule implementations can retrieve the request URL via
	// resp.Request.URL.
	Apply(physurl *url.URL, resp *http.Response) (*url.URL, error)
}

// AppendExtDotUnixTime generates the validity URL by appending the provided
// extension ext and the last modified time to the physical URL. For example,
//
//     https://example.com/index.html
//
// may yield the following validity URL:
//
//     https://example.com/index.html.validity.1561984496
//
// ext usually starts with a dot ("."); AppendExtDotUnixTime does not insert
// the one automatically. ext is therefore ".validity" rather than "validity"
// in the example above.
//
// The last modified time is represented by a UNIX timestamp, and taken from
// the Last-Modified header of the response. When Last-Modified is missing or
// has an unparsable value, AppendExtDotUnixTime uses the now argument instead.
//
// AppendExtDotUnixTime does not expect the physical URL to contain Query or
// Fragment. They will be stripped off from the validity URL.
//
// Apply returns an error when physurl looks like a directory, e.g. the path
// ends with a slash.
func AppendExtDotUnixTime(ext string, now time.Time) ValidityURLRule {
	return &appendExtDotUnixTime{ext, now}
}

type appendExtDotUnixTime struct {
	ext string
	now time.Time
}

func (rule *appendExtDotUnixTime) Apply(physurl *url.URL, resp *http.Response) (*url.URL, error) {
	// We do not care whether physurl is normalized or not: we can append
	// the extension as long as it has a filename.
	if urlutil.IsDir(physurl) {
		return nil, fmt.Errorf("%q looks like a directory", physurl)
	}
	date := rule.parseDate(resp)
	newPath := fmt.Sprintf("%s%s.%d", physurl.Path, rule.ext, date.Unix())
	return physurl.ResolveReference(&url.URL{Path: newPath}), nil
}

func (rule *appendExtDotUnixTime) parseDate(resp *http.Response) time.Time {
	date := resp.Header.Get("Last-Modified")
	if date == "" {
		return rule.now
	}
	parsed, err := httpheader.ParseDate(date)
	if err != nil {
		log.Printf("warning: failed to parse the header %q: %v", date, err)
		return rule.now
	}
	return parsed
}
