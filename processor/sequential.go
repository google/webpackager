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

package processor

import (
	"github.com/google/webpackager/exchange"
)

// SequentialProcessor consists of a series of subprocessors.
type SequentialProcessor []Processor

// Process invokes all subprocessors in order.
//
// Process fails immediately when some subprocessor returns an error.
// The subsequent subprocessors will not run in such case.
func (sp SequentialProcessor) Process(resp *exchange.Response) error {
	for _, p := range sp {
		if err := p.Process(resp); err != nil {
			return err
		}
	}
	return nil
}
