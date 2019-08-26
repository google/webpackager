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

package processor_test

import (
	"reflect"
	"testing"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor"
)

func TestMultiplexedProcessor(t *testing.T) {
	mp := processor.MultiplexedProcessor{
		"text/html":  newTestingProcessor("html"),
		"image/gif":  newTestingProcessor("gif"),
		"image/jpeg": newTestingProcessor("jpeg"),
		"image/png":  newTestingProcessor("png"),
	}

	tests := []struct {
		ctype string
		want  []string
	}{
		{
			ctype: "text/html",
			want:  []string{"html"},
		},
		{
			ctype: "Text/HTML",
			want:  []string{"html"},
		},
		{
			ctype: "text/html; charset=utf-8",
			want:  []string{"html"},
		},
		{
			ctype: "text/html; charset=(@_@)",
			want:  []string{"html"},
		},
		{
			ctype: "image/jpeg",
			want:  []string{"jpeg"},
		},
		{
			ctype: "",
			want:  nil,
		},
		{
			ctype: "text/plain",
			want:  nil,
		},
		{
			ctype: "text/plain; charset=utf-8",
			want:  nil,
		},
		{
			ctype: "text/html :)", // Invalid
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.ctype, func(t *testing.T) {
			resp := exchangetest.MakeEmptyResponse("https://dummy.test/")
			if test.ctype != "" {
				resp.Header.Set("Content-Type", test.ctype)
			}
			if err := mp.Process(resp); err != nil {
				t.Errorf("got error(%q), want success", err)
			}
			if got := resp.Header["X-Testing"]; !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
