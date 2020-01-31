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

package customflag_test

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/internal/customflag"
)

func TestMultiString(t *testing.T) {
	fs := flag.NewFlagSet("multistring_test", flag.ContinueOnError)
	var foo []string
	fs.Var(customflag.NewMultiStringValue(&foo), "foo", "Test flag #1.")
	var bar []string
	fs.Var(customflag.NewMultiStringValue(&bar), "bar", "Test flag #2.")

	fs.Parse([]string{
		"--foo=a", "--bar=b", "--foo=c", "--foo=d", "--bar=e",
	})

	if diff := cmp.Diff([]string{"a", "c", "d"}, foo); diff != "" {
		t.Errorf("foo mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"b", "e"}, bar); diff != "" {
		t.Errorf("bar mismatch (-want +got):\n%s", diff)
	}
}
