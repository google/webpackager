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

package filewrite

import (
	"github.com/google/webpackager/resource/cache"
)

// Config holds the parameters to NewFileWriteCache.
type Config struct {
	// BaseCache specifies the underlying ResourceCache.
	BaseCache cache.ResourceCache

	// ExchangeMapping specifies the rule to determine the location of signed
	// exchange files. nil is equivalent to MapToDevNull.
	ExchangeMapping MappingRule

	// ValidityMapping is currently unused.
	ValidityMapping MappingRule
}
