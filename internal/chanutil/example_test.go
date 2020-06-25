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

package chanutil_test

import (
	"fmt"
	"sync"

	"github.com/google/webpackager/internal/chanutil"
)

func ExampleKiller() {
	k := chanutil.NewKiller()

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func(i int, k *chanutil.Killer) {
			defer wg.Done()
			<-k.C // Wait for the kill signal.
			fmt.Printf("Worker %v is killed\n", i)
		}(i, k)
	}

	k.Kill()

	wg.Wait()
	// Unordered output:
	// Worker 0 is killed
	// Worker 1 is killed
	// Worker 2 is killed
}
