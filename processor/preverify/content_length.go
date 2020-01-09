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

package preverify

import (
	"fmt"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
)

// MaxContentLength requires the content (the response body) to be not
// larger then limit.
func MaxContentLength(limit int) processor.Processor {
	return &maxContentLength{limit}
}

type maxContentLength struct {
	limit int
}

func (mcl *maxContentLength) Process(resp *exchange.Response) error {
	if len(resp.Payload) > mcl.limit {
		return fmt.Errorf("oversized content (%d bytes; limit: %d bytes)",
			len(resp.Payload), mcl.limit)
	}
	return nil
}
