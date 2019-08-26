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
	"errors"
	"reflect"
	"testing"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor"
)

func TestSequentialProcessor(t *testing.T) {
	errDummy := errors.New("errDummy")

	tests := []struct {
		name string
		proc processor.SequentialProcessor
		err  error
		want []string
	}{
		{
			name: "NoError",
			proc: processor.SequentialProcessor{
				newTestingProcessor("foo"),
				newTestingProcessor("bar"),
				newTestingProcessor("baz"),
			},
			err:  nil,
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "Error",
			proc: processor.SequentialProcessor{
				newTestingProcessor("foo"),
				newTestingProcessor("bar"),
				newFailingProcessor(errDummy),
				newTestingProcessor("baz"),
			},
			err:  errDummy,
			want: []string{"foo", "bar"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := exchangetest.MakeEmptyResponse("https://dummy.test/")

			if err := test.proc.Process(resp); err != test.err {
				t.Errorf("got %v, want %v", err, test.err)
			}
			if got := resp.Header["X-Testing"]; !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
