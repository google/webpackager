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

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/internal/urlutil"
)

// AppendExtDotLastModified generates the validity URL by appending ext
// and the resource's last modified time. For example:
//
//     https://example.com/index.html
//
// would receive a validity URL that looks like:
//
//     https://example.com/index.html.validity.1561984496
//
// ext usually starts with a dot ("."). AppendExtDotExchangeDate does not
// insert it automatically. ext is thus ".validity" rather than "validity"
// in the example above.
//
// The last modified time is taken from the Last-Modified header field in
// the HTTP response and represented in UNIX time. If the Last-Modified
// is missing or unparsable, AppendExtDotLastModified uses the date value
// of the signed exchange signature (vp.Date).
//
// The AppendExtDotLastModified rule does not support physurl that looks
// like a directory (e.g. has Path ending with a slash). Apply returns an
// error for such physurl. Note you can use urlrewrite.IndexRule to ensure
// physurl to always have a filename.
//
// The AppendExtDotLastModified rule ignores Query and Fragment in physurl.
// The validity URLs will always have empty Query and Fragment.
func AppendExtDotLastModified(ext string) URLRule {
	return &appendExtDotLastModified{ext}
}

type appendExtDotLastModified struct {
	ext string
}

func (rule *appendExtDotLastModified) Apply(physurl *url.URL, resp *exchange.Response, vp exchange.ValidPeriod) (*url.URL, error) {
	date := resp.Header.Get("Last-Modified")
	if date == "" {
		return toValidityURL(physurl, rule.ext, vp.Date())
	}
	parsed, err := http.ParseTime(date)
	if err != nil {
		log.Printf("warning: failed to parse the header %q: %v", date, err)
		return toValidityURL(physurl, rule.ext, vp.Date())
	}
	return toValidityURL(physurl, rule.ext, parsed)
}

// AppendExtDotExchangeDate is like AppendExtDotLastModified but always
// uses vp.Date instead of the last modified time.
func AppendExtDotExchangeDate(ext string) URLRule {
	return &appendExtDotExchangeDate{ext}
}

type appendExtDotExchangeDate struct {
	ext string
}

func (rule *appendExtDotExchangeDate) Apply(physurl *url.URL, resp *exchange.Response, vp exchange.ValidPeriod) (*url.URL, error) {
	return toValidityURL(physurl, rule.ext, vp.Date())
}

func toValidityURL(physurl *url.URL, ext string, date time.Time) (*url.URL, error) {
	// We do not care whether physurl is normalized or not: we can append
	// the extension as long as it has a filename.
	if urlutil.IsDir(physurl) {
		return nil, fmt.Errorf("%q looks like a directory", physurl)
	}
	newPath := fmt.Sprintf("%s%s.%d", physurl.Path, ext, date.Unix())
	// This ResolveReference drops Query and Fragment from the resulting URL.
	return physurl.ResolveReference(&url.URL{Path: newPath}), nil
}
