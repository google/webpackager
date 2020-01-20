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

/*
Package htmltask implements some optimization logic for HTML documents.

The logic is implemented on a common interface named HTMLTask. Each HTMLTask
implementation has a single clear focus, such as "add the preload links for
stylesheets used in the HTML document." HTMLTasks are passed collectively to
htmlproc.NewHTMLProcessor to define its processing logic.
*/
package htmltask

import (
	"github.com/google/webpackager/processor/htmlproc/htmldoc"
)

// HTMLTask manipulates HTMLResponse to optimize the page loading.
type HTMLTask interface {
	Run(resp *htmldoc.HTMLResponse) error
}

// ConservativeTaskSet is the set of HTMLTasks used in the default config.
// It consists only of HTMLTasks that almost always work well.
var ConservativeTaskSet = []HTMLTask{
	ExtractSubContentTypes(),
	ExtractPreloadTags(),
}

// AggressiveTaskSet gets as many resources preloaded as Web Packager can.
// It includes HTMLTasks that might make negative effect in some cases.
var AggressiveTaskSet = []HTMLTask{
	ExtractSubContentTypes(),
	ExtractPreloadTags(),
	PreloadStylesheets(),
	InsecurePreloadScripts(),
}
