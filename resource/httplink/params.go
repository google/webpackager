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

package httplink

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// A few common parameter names for use with LinkParams.
const (
	ParamRel         = "rel"
	ParamAs          = "as"
	ParamCrossOrigin = "crossorigin"
	ParamMedia       = "media"
	ParamType        = "type"
)

// Special parameter values recognized by LinkParams.
const (
	// Value(s) for the "rel" parameter.
	RelPreload = "preload"

	// Value(s) for the "crossorigin" parameter.
	CrossOriginAnonymous = "anonymous"
)

// LinkParams represents the parameters of a Web Linking.
type LinkParams map[string]string

// Get returns the value of the parameter specified by key. key gets lowered,
// thus is case-insensitive.
func (p LinkParams) Get(key string) string {
	return p[strings.ToLower(key)]
}

// Set changes the value of the parameter specified by key. key gets lowered,
// thus is case-insensitive. Set also normalizes the provided value for some
// parameters, e.g. removes extra spaces for the rel parameter. To get around
// the normalization, access the map entry directly.
func (p LinkParams) Set(key, val string) {
	key = strings.ToLower(key)
	p[key] = normalizeValue(key, val)
}

// Clone returns a deep copy of the LinkParams p.
func (p LinkParams) Clone() LinkParams {
	q := make(LinkParams, len(p))
	for k, v := range p {
		q[k] = v
	}
	return q
}

// Equal reports whether p and q contain the same set of key-value pairs.
func (p LinkParams) Equal(q LinkParams) bool {
	return reflect.DeepEqual(p, q)
}

func (p LinkParams) write(w io.Writer) {
	keys := make([]string, 0, len(p))
	for key := range p {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[j] == ParamRel {
			return false
		}
		if keys[i] == ParamRel {
			return true
		}
		return keys[i] < keys[j]
	})
	for _, key := range keys {
		if shouldElideValue(key, p[key]) {
			fmt.Fprintf(w, ";%s", key)
		} else {
			fmt.Fprintf(w, ";%s=%q", key, p[key])
		}
	}
}

func normalizeValue(key, val string) string {
	switch key {
	case ParamRel:
		// [RFC 8288] requires the relation types to be compared character
		// by character in a case-insensitive fashion, whether they are
		// registered (well-known) or external (represented by URIs). Note
		// also [RFC 8288] recommends URIs to be all lowercase.
		//
		// [RFC 8288]: https://tools.ietf.org/html/rfc8288
		vals := strings.Fields(val)
		for i := range vals {
			vals[i] = strings.ToLower(vals[i])
		}
		return strings.Join(vals, " ")

	case ParamAs, ParamType:
		return strings.ToLower(val)

	case ParamCrossOrigin:
		if val == "" {
			return CrossOriginAnonymous
		}
		return strings.ToLower(val)

	default:
		return val // Do not normalize unknown parameters.
	}
}

func shouldElideValue(key, val string) bool {
	// Elide the value from `crossorigin="anonymous"` since the bare
	// `crossorigin` is presumably more popular.
	return key == ParamCrossOrigin && val == CrossOriginAnonymous
}
