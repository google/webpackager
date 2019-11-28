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

package fetch

import (
	"net/http"
)

func clone(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// RequestTweaker mutates an http.Request to set the refer(r)er, add custom
// HTTP headers, and so on.
type RequestTweaker interface {
	// Tweak mutates req. parent is the request which spawned req. If req is
	// a request to fetch a subresource, parent is typically the one used to
	// fetch the main resource. parent can be nil (e.g. when req is a request
	// to fetch a main resource) and should not be mutated.
	Tweak(req *http.Request, parent *http.Request) error
}

// DefaultRequestTweaker is a RequestTweaker used by default.
var DefaultRequestTweaker RequestTweaker = SetReferer()

// RequestTweakerSequence consists of a series of RequestTweakers.
type RequestTweakerSequence []RequestTweaker

// Tweak invokes all tweakers in order. It fails immediately when some tweaker
// returns an error, in which case the subsequent tweakers will not run.
func (seq RequestTweakerSequence) Tweak(req, parent *http.Request) error {
	for _, tweaker := range seq {
		if err := tweaker.Tweak(req, parent); err != nil {
			return err
		}
	}
	return nil
}

// SetReferer sets the Referer HTTP header with the parent request URL.
func SetReferer() RequestTweaker {
	return &setReferer{}
}

type setReferer struct{}

func (*setReferer) Tweak(req, parent *http.Request) error {
	if parent != nil {
		req.Header.Set("Referer", parent.URL.String())
	}
	return nil
}

// CopyParentHeaders copies the header fields of the provided keys from the
// parent request. When the tweaked request already has those header fields,
// their values will be overwritten by the values from the parent request.
// CopyParentHeaders, however, only mutates the header fields present in the
// parent request and keeps all other header fields untouched.
//
// keys are case insensitive: they are canonicalized by http.CanonicalHeaderKey.
func CopyParentHeaders(keys []string) RequestTweaker {
	canonicalKeys := make([]string, len(keys))
	for i, key := range keys {
		canonicalKeys[i] = http.CanonicalHeaderKey(key)
	}
	return &copyParentHeaders{canonicalKeys}
}

type copyParentHeaders struct {
	keys []string
}

func (cph *copyParentHeaders) Tweak(req, parent *http.Request) error {
	if parent != nil {
		for _, key := range cph.keys {
			if v, ok := parent.Header[key]; ok {
				req.Header[key] = clone(v)
			}
		}
	}
	return nil
}

// SetCustomHeaders populates the provided HTTP header fields to the request.
// When the request already has the same header fields, their values will be
// overwritten with the values from the provided http.Header.
func SetCustomHeaders(header http.Header) RequestTweaker {
	return &setCustomHeaders{header}
}

type setCustomHeaders struct {
	header http.Header
}

func (sch *setCustomHeaders) Tweak(req, parent *http.Request) error {
	for k, v := range sch.header {
		req.Header[k] = clone(v)
	}
	return nil
}
