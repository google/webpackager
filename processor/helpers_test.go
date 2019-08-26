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
	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
)

// newTestingProcessor returns a Processor that adds "X-Testing: <value>"
// to the response header.
func newTestingProcessor(value string) processor.Processor {
	return &testingProcessor{value}
}

type testingProcessor struct {
	value string
}

func (p *testingProcessor) Process(resp *exchange.Response) error {
	resp.Header.Add("X-Testing", p.value)
	return nil
}

// newFailingProcessor returns a Processor that always fails with err.
func newFailingProcessor(err error) processor.Processor {
	return &failingProcessor{err}
}

type failingProcessor struct {
	err error
}

func (p *failingProcessor) Process(resp *exchange.Response) error {
	return p.err
}
