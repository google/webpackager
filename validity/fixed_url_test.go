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

package validity_test

import (
	"net/url"
	"testing"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/validity"
)

func TestFixedURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		rule validity.URLRule
		want string
	}{
		{
			name: "AbsoluteURL",
			url:  "https://example.com/index.html",
			rule: validity.FixedURL(
				urlutil.MustParse("https://example.com/empty.validity"),
			),
			want: "https://example.com/empty.validity",
		},
		{
			name: "RelativeURL",
			url:  "https://example.com/index.html",
			rule: validity.FixedURL(
				urlutil.MustParse("/empty.validity"),
			),
			want: "https://example.com/empty.validity",
		},
	}

	for _, test := range tests {
		arg, err := url.Parse(test.url)
		if err != nil {
			t.Fatal(err)
		}
		resp := exchangetest.MakeEmptyResponse(test.url)
		// FixedURL does not use the ValidPeriod.
		got, err := test.rule.Apply(arg, resp, exchange.ValidPeriod{})
		if err != nil {
			t.Fatalf("got error(%q), want success", err)
		}
		if got.String() != test.want {
			t.Errorf("got %q, want %q", got, test.want)
		}
	}
}
