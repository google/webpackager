// Copyright 2020 Google LLC
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

package certmanager

import (
	"time"

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

// DummyOCSPRespSource always returns certchain.DummyOCSPResponse.
var DummyOCSPRespSource OCSPRespSource = &dummyOCSPRespSource{}

type dummyOCSPRespSource struct{}

func (*dummyOCSPRespSource) Fetch(chain *certchain.RawChain, now func() time.Time) (ocspResp *certchain.OCSPResponse, nextRun futureevent.Event, err error) {
	return certchain.DummyOCSPResponse, futureevent.NeverOccurs(), nil
}
