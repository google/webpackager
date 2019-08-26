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
	"net/http"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
)

// RequireStatusOK ensures the response to have the status code 200 (OK).
var RequireStatusOK processor.Processor = &requireStatusOK{}

type requireStatusOK struct{}

func (*requireStatusOK) Process(resp *exchange.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("responded with status %d", resp.StatusCode)
	}
	return nil
}
