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

package htmlproc_test

import (
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"github.com/google/webpackager/processor/htmlproc/htmltask"
)

type HTMLTaskFunc func(*htmldoc.HTMLResponse) error

func AsHTMLTask(run HTMLTaskFunc) htmltask.HTMLTask {
	return &asHTMLTask{run}
}

type asHTMLTask struct {
	run HTMLTaskFunc
}

func (task *asHTMLTask) Run(resp *htmldoc.HTMLResponse) error {
	return task.run(resp)
}
