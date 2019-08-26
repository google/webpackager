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

// Package htmlproc implements a Processor to process HTML documents.
package htmlproc

import (
	"bytes"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
	"github.com/google/webpackager/processor/htmlproc/htmltask"
	"golang.org/x/net/html"
)

// Config holds parameters to NewHTMLProcessor.
type Config struct {
	// TaskSet specifies the sequence of HTMLTasks to run.
	TaskSet []htmltask.HTMLTask

	// ModifyHTML indicates whether the processor can modify HTML documents.
	//
	// When ModifyHTML is true, the processor always rewrites the payload
	// with HTML reconstructed from the parse tree. The response thus always
	// contains a well-formed HTML after processing.
	//
	// Some HTMLTasks have an effect only when ModifyHTML is true.
	ModifyHTML bool
}

// NewHTMLProcessor creates and initializes a new Processor to process HTML
// documents. The Processor turns provided exchange.Response into
// htmldoc.HTMLResponse then runs the specified htmltask.HTMLTasks one by one.
// The Processor fails immediately when some HTMLTask encounters an error.
func NewHTMLProcessor(config Config) processor.Processor {
	return &htmlProcessor{config}
}

type htmlProcessor struct {
	Config
}

func (hp *htmlProcessor) Process(resp *exchange.Response) error {
	htmlResp, err := htmldoc.NewHTMLResponse(resp)
	if err != nil {
		return err
	}

	for _, task := range hp.TaskSet {
		if err := task.Run(htmlResp); err != nil {
			return err
		}
	}

	if hp.ModifyHTML {
		var payload bytes.Buffer
		if err := html.Render(&payload, htmlResp.Doc.Root); err != nil {
			return err
		}
		htmlResp.Payload = payload.Bytes()
	}

	return nil
}
