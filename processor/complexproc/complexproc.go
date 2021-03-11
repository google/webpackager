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

package complexproc

import (
	"github.com/google/webpackager/processor"
	"github.com/google/webpackager/processor/commonproc"
	"github.com/google/webpackager/processor/htmlproc"
	"github.com/google/webpackager/processor/preverify"
)

// DefaultProcessor is the processor used by webpackager.Packager by default.
var DefaultProcessor = NewComprehensiveProcessor(Config{})

// Config customizes NewComprehensiveProcessor.
type Config struct {
	// Preverify is passed to preverify.CheckPrerequisites.
	Preverify preverify.Config

	// HTML is passed to htmlproc.NewHTMLProcessor.
	HTML htmlproc.Config

	// CustomMainProcessors is a map from media types to main processors.
	//
	// CustomMainProcessors takes the precedence over the default map.
	// When both maps contain the same keys, ComprehensiveProcessor runs
	// the processors from CustomMainProcessors instead of the default ones.
	//
	// CustomMainProcessors can have entries with a nil value. They disable
	// the default processors for the given media types.
	CustomMainProcessors processor.MultiplexedProcessor

	// CustomPreprocessors are run before the main processor.
	CustomPreprocessors processor.SequentialProcessor

	// CustomPostprocessors are run after the main processor.
	CustomPostprocessors processor.SequentialProcessor
}

// These processors are always included in ComprehensiveProcessors.
var (
	// EssentialPreprocessors contain always-run preprocessors.
	EssentialPreprocessors = processor.SequentialProcessor{
		commonproc.ExtractPreloadHeaders,
	}
	// EssentialPostprocessors contain always-run postprocessors.
	EssentialPostprocessors = processor.SequentialProcessor{
		commonproc.ContentTypeProcessor,
		commonproc.RemoveUncachedHeaders,
	}
)

// NewComprehensiveProcessor creates and initializes a new processor based
// on the provided Config.
func NewComprehensiveProcessor(config Config) processor.Processor {
	// TODO(yuizumi): Maybe flatten these processors.
	return processor.SequentialProcessor{
		preverify.CheckPrerequisites(config.Preverify),
		EssentialPreprocessors,
		config.CustomPreprocessors,
		newMainProcessor(config),
		EssentialPostprocessors,
		config.CustomPostprocessors,
	}
}

func newMainProcessor(config Config) processor.Processor {
	// TODO(yuizumi): Add processors for other types (e.g. images).
	html := htmlproc.NewHTMLProcessor(config.HTML)
	mp := processor.MultiplexedProcessor{
		"text/html":             html,
		"application/xhtml+xml": html,
	}
	for k, v := range config.CustomMainProcessors {
		if v == nil {
			delete(mp, k)
		} else {
			mp[k] = v
		}
	}
	return mp
}
