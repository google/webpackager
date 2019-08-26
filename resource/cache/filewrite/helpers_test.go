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

package filewrite_test

import (
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/cache/filewrite"
)

// FixedMappingRule returns a MappingRule mapping anything to path.
func FixedMappingRule(path string) filewrite.MappingRule {
	return &fixedMappingRule{path}
}

type fixedMappingRule struct {
	path string
}

func (rule *fixedMappingRule) Map(r *resource.Resource) (string, error) {
	return rule.path, nil
}

// ErrorMappingRule returns a MappingRule failing with err all the time.
func ErrorMappingRule(err error) filewrite.MappingRule {
	return &errorMappingRule{err}
}

type errorMappingRule struct {
	err error
}

func (rule *errorMappingRule) Map(r *resource.Resource) (string, error) {
	return "", rule.err
}
