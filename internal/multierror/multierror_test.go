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

package multierror_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/webpackager/internal/multierror"
)

func TestMultiError(t *testing.T) {
	var me multierror.MultiError

	errFoo := errors.New("test: foo")
	me.Add(errFoo)
	errBar := errors.New("test: bar")
	me.Add(errBar)
	errBaz := errors.New("test: baz")
	me.Add(errBaz)

	wantErrors := []error{errFoo, errBar, errBaz}
	if !reflect.DeepEqual(me.Errors, wantErrors) {
		t.Errorf("me.Errors = %q, want %q", me.Errors, wantErrors)
	}

	wantText := "test: foo; test: bar; test: baz"
	if got := me.Error(); got != wantText {
		t.Errorf("me.Error() = %q, want %q", got, wantText)
	}
}
