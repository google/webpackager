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

package vprule_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/exchange/vprule"
)

func TestPreContentType(t *testing.T) {
	rule7Day := vprule.FixedLifetime(7 * 24 * time.Hour)
	rule2Day := vprule.FixedLifetime(2 * 24 * time.Hour)
	rule1Day := vprule.FixedLifetime(1 * 24 * time.Hour)

	tests := []struct {
		name  string
		resp  string
		extra http.Header
		rule  vprule.Rule
		date  time.Time
		want  exchange.ValidPeriod
	}{
		{
			name: "ContentType_NotMatching",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Content-Length: 35\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><p>Hello, world!</p>",
			),
			rule: vprule.PerContentType(
				map[string]vprule.Rule{
					"text/css":               rule2Day,
					"application/javascript": rule1Day,
					"text/javascript":        rule1Day,
				},
				rule7Day,
			),
			date: time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
			want: exchange.NewValidPeriod(
				time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
				time.Date(2020, time.January, 22, 19, 30, 0, 0, time.UTC)),
		},
		{
			name: "ContentType_Matching",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Content-Length: 35\r\n",
				"Content-Type: text/css\r\n",
				"\r\n",
				"body { background-color: #abcdef; }",
			),
			rule: vprule.PerContentType(
				map[string]vprule.Rule{
					"text/css":               rule2Day,
					"application/javascript": rule1Day,
					"text/javascript":        rule1Day,
				},
				rule7Day,
			),
			date: time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
			want: exchange.NewValidPeriod(
				time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
				time.Date(2020, time.January, 17, 19, 30, 0, 0, time.UTC)),
		},
		{
			name: "SubContentType_NotMatching",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Content-Length: 38\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><math><mn>42</mn></math>",
			),
			extra: http.Header{
				exchange.SubContentType: []string{"application/mathml+xml"},
			},
			rule: vprule.PerContentType(
				map[string]vprule.Rule{
					"text/css":               rule2Day,
					"application/javascript": rule1Day,
					"text/javascript":        rule1Day,
				},
				rule7Day,
			),
			date: time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
			want: exchange.NewValidPeriod(
				time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
				time.Date(2020, time.January, 22, 19, 30, 0, 0, time.UTC)),
		},
		{
			name: "SubContentType_Matching",
			resp: fmt.Sprint(
				"HTTP/1.1 200 OK\r\n",
				"Content-Length: 39\r\n",
				"Content-Type: text/html; charset=utf-8\r\n",
				"\r\n",
				"<!doctype html><script>/*...*/</script>",
			),
			extra: http.Header{
				exchange.SubContentType: []string{"application/javascript"},
			},
			rule: vprule.PerContentType(
				map[string]vprule.Rule{
					"text/css":               rule2Day,
					"application/javascript": rule1Day,
					"text/javascript":        rule1Day,
				},
				rule7Day,
			),
			date: time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
			want: exchange.NewValidPeriod(
				time.Date(2020, time.January, 15, 19, 30, 0, 0, time.UTC),
				time.Date(2020, time.January, 16, 19, 30, 0, 0, time.UTC)),
		},
	}

	url := "https://example.com/"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeResponse(url, test.resp)
			resp.ExtraData = test.extra

			if got := test.rule.Get(resp, test.date); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
