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

// Package defaultproc provides DefaultProcessor, a processor that can be
// used out of the box. It is defined as a composite of multiple processors,
// which can be reused to construct slightly modified processors.
package defaultproc

import (
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/commonproc"
	"github.com/google/webpackager/processor/htmlproc"
	"github.com/google/webpackager/processor/preverify"
)

// DefaultProcessor is the processor used by webpackager.Packager by default.
var DefaultProcessor = processor.SequentialProcessor{
	Preprocessors,
	NewMainProcessor(Config{}),
	Postprocessors,
}

// Config holds the settings passed to NewMainProcessor.
type Config struct {
	HTML htmlproc.Config
}

// Here are the parameters to DefaultProcessor.
var (
	// Preprocessors specifies the processors to run before the main processor.
	Preprocessors = processor.SequentialProcessor{
		preverify.CheckPrerequisites,
		commonproc.ExtractPreloadHeaders,
	}
	// Postprocessors specifies the processors to run after the main processor.
	Postprocessors = processor.SequentialProcessor{
		commonproc.ApplySameOriginPolicy,
		commonproc.ContentTypeProcessor,
		commonproc.RemoveUncachedHeaders,
	}
)

// NewMainProcessor creates and initializes a MultiplexedProcessor based on
// the provided Config. It just contains HTMLProcessor applied to text/html
// (HTML) and application/xhtml+xml (XHTML), at this moment.
func NewMainProcessor(config Config) processor.MultiplexedProcessor {
	// TODO(yuizumi): Add processors for other types (e.g. images).
	html := htmlproc.NewHTMLProcessor(config.HTML)
	return processor.MultiplexedProcessor{
		"text/html":             html,
		"application/xhtml+xml": html,
	}
}
