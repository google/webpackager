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

package commonproc

import (
	"log"
	"net/http"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/processor"
)

// ContentTypeProcessor adds the "X-Content-Type-Options: nosniff" header,
// which prevents the browsers to "auto-correct" the content type.
// ContentTypeProcessor also adds the Content-Type header when it is absent,
// emitting a warning log.
var ContentTypeProcessor processor.Processor = &contentTypeProcessor{}

type contentTypeProcessor struct{}

func (*contentTypeProcessor) Process(resp *exchange.Response) error {
	resp.Header.Set("X-Content-Type-Options", "nosniff")

	if resp.Header.Get("Content-Type") == "" {
		ctype := http.DetectContentType(resp.Payload)
		log.Printf("warning: %s is missing Content-Type; set to %q.",
			resp.Request.URL, ctype)
		resp.Header.Set("Content-Type", ctype)
	}

	return nil
}
