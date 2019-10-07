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

package exchange_test

import (
	"fmt"
	"time"

	"github.com/google/webpackager/exchange"
)

func ExampleValidPeriod_Contains() {
	date := time.Date(2019, time.October, 1, 10, 30, 0, 0, time.UTC)
	expires := time.Date(2019, time.October, 2, 10, 30, 0, 0, time.UTC)

	vp := exchange.NewValidPeriod(date, expires)

	fmt.Println(vp.Contains(time.Date(2019, time.October, 1, 10, 29, 0, 0, time.UTC)))
	fmt.Println(vp.Contains(time.Date(2019, time.October, 1, 10, 30, 0, 0, time.UTC)))
	fmt.Println(vp.Contains(time.Date(2019, time.October, 1, 22, 30, 0, 0, time.UTC)))
	fmt.Println(vp.Contains(time.Date(2019, time.October, 2, 10, 30, 0, 0, time.UTC)))
	fmt.Println(vp.Contains(time.Date(2019, time.October, 2, 10, 31, 0, 0, time.UTC)))
	// Output:
	// false
	// true
	// true
	// true
	// false
}
