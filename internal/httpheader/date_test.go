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

package httpheader_test

import (
	"testing"
	"time"

	"github.com/google/webpackager/internal/httpheader"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want time.Time
	}{
		{
			name: "IMF-fixdate",
			arg:  "Sun, 06 Nov 1994 08:49:37 GMT",
			want: time.Date(1994, time.November, 6, 8, 49, 37, 0, time.UTC),
		},
		{
			name: "IMF-fixdate_incomplete",
			arg:  "6 Nov 1994 08:49 +0000",
			want: time.Date(1994, time.November, 6, 8, 49, 0, 0, time.UTC),
		},
		{
			name: "rfc850-date",
			arg:  "Sunday, 06-Nov-94 08:49:37 GMT",
			want: time.Date(1994, time.November, 6, 8, 49, 37, 0, time.UTC),
		},
		{
			name: "asctime-date",
			arg:  "Sun Nov  6 08:49:37 1994",
			want: time.Date(1994, time.November, 6, 8, 49, 37, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := httpheader.ParseDate(test.arg)
			if err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if !got.Equal(test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
