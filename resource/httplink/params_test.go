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

package httplink_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/resource/httplink"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name   string
		params httplink.LinkParams
		key    string
		want   string
	}{
		{
			name:   "Existent",
			params: httplink.LinkParams{"rel": "preload", "type": "text/css"},
			key:    "type",
			want:   "text/css",
		},
		{
			name:   "CaseIgnored",
			params: httplink.LinkParams{"rel": "preload", "type": "text/css"},
			key:    "TyPe",
			want:   "text/css",
		},
		{
			name:   "Nonexistent",
			params: httplink.LinkParams{"rel": "preload"},
			key:    "type",
			want:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.params.Get(test.key)
			if got != test.want {
				t.Errorf("params.Get(%q) = %q, want %q", test.key, got, test.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		from httplink.LinkParams
		key  string
		val  string
		want httplink.LinkParams
	}{
		{
			name: "Insert",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "as",
			val: "style",
			want: httplink.LinkParams{
				"rel": "preload", "as": "style",
			},
		},
		{
			name: "Overwrite",
			from: httplink.LinkParams{
				"rel": "preload", "as": "image",
			},
			key: "as",
			val: "style",
			want: httplink.LinkParams{
				"rel": "preload", "as": "style",
			},
		},
		{
			name: "Normalize_Rel",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "rel",
			val: " StyleSheet  ALTeRNATe",
			want: httplink.LinkParams{
				"rel": "stylesheet alternate",
			},
		},
		{
			name: "Normalize_As",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "as",
			val: "STyLe",
			want: httplink.LinkParams{
				"rel": "preload", "as": "style",
			},
		},
		{
			name: "Normalize_CrossOrigin_Anonymous",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "crossorigin",
			val: "",
			want: httplink.LinkParams{
				"rel": "preload", "crossorigin": "anonymous",
			},
		},
		{
			name: "Normalize_CrossOrigin_Lowercase",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "crossorigin",
			val: "USER-CREDENTIALS",
			want: httplink.LinkParams{
				"rel": "preload", "crossorigin": "user-credentials",
			},
		},
		{
			name: "Normalize_Type",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "type",
			val: "text/CSS",
			want: httplink.LinkParams{
				"rel": "preload", "type": "text/css",
			},
		},
		{
			name: "Normalize_KeyVal",
			from: httplink.LinkParams{
				"rel": "preload",
			},
			key: "TyPe",
			val: "text/CSS",
			want: httplink.LinkParams{
				"rel": "preload", "type": "text/css",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			params := test.from.Clone()
			params.Set(test.key, test.val)

			// Convert to map[string]string for a better message.
			var want map[string]string = test.want
			var got map[string]string = params
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("params mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
