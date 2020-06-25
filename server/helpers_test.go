// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server_test

import (
	"time"

	"github.com/google/webpackager/certchain"
	"github.com/google/webpackager/certchain/certmanager/futureevent"
)

type stubRawChainSource struct {
	data *certchain.RawChain
}

func (s *stubRawChainSource) Fetch(*certchain.RawChain, func() time.Time) (*certchain.RawChain, futureevent.Event, error) {
	return s.data, futureevent.NeverOccurs(), nil
}

type stubOCSPRespSource struct {
	data *certchain.OCSPResponse
}

func (s *stubOCSPRespSource) Fetch(*certchain.RawChain, func() time.Time) (*certchain.OCSPResponse, futureevent.Event, error) {
	return s.data, futureevent.NeverOccurs(), nil
}
