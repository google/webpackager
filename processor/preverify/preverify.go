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

// Package preverify implements processors to verify that HTTP responses
// can be distributed as signed exchanges. These processors do not mutate
// the provided exchange.Response: they just inspect it and report an error
// when it does not meet the criteria.
package preverify

import (
	"github.com/google/webpackager/processor"
)

// CheckPrerequisites is a composite of all prerequisite checkers.
var CheckPrerequisites = processor.SequentialProcessor{
	RequireStatusOK,
}
