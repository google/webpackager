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

package defaultproc_test

import (
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/defaultproc"
	"github.com/google/webpackager/processor/htmlproc"
	"github.com/google/webpackager/processor/htmlproc/htmltask"
)

// This example constructs a new Processor that runs a custom HTMLTask and
// behaves otherwise the same as DefaultProcessor.
func Example_customize() processor.Processor {
	// Instantiate your custom HTMLTask.
	yourTask := NewCustomHTMLTask()

	// Have your HTMLTask run in the HTMLProcessor.
	config := defaultproc.Config{
		HTML: htmlproc.Config{
			TaskSet: append(htmltask.DefaultTaskSet, yourTask),
		},
	}

	// Create the composite with the config above.
	return processor.SequentialProcessor{
		defaultproc.Preprocessors,
		defaultproc.NewMainProcessor(config),
		defaultproc.Postprocessors,
	}
}
