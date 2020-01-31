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

package preload_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/preload"
)

func newResource(rawurl string) *resource.Resource {
	return resource.NewResource(urlutil.MustParse(rawurl))
}

func TestNewPlainPreload(t *testing.T) {
	r := newResource("https://example.com/test.css")
	p := preload.NewPlainPreload(r, preload.AsStyle)

	t.Run("Resources", func(t *testing.T) {
		if diff := cmp.Diff([]*resource.Resource{r}, p.Resources()); diff != "" {
			t.Errorf("p.Resources() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("Header", func(t *testing.T) {
		want := `<https://example.com/test.css>;rel="preload";as="style"`
		if got := p.Header(); got != want {
			t.Errorf("p.Header() = %q, want %q", got, want)
		}
	})
}

func TestCompletePlainPreload(t *testing.T) {
	r := newResource("https://example.com/font.woff2")

	p := &preload.PlainPreload{
		Resource:    r,
		As:          preload.AsFont,
		CrossOrigin: true,
		Media:       "screen",
		Type:        "font/woff2",
	}

	t.Run("Resources", func(t *testing.T) {
		if diff := cmp.Diff([]*resource.Resource{r}, p.Resources()); diff != "" {
			t.Errorf("p.Resources() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("Header", func(t *testing.T) {
		want := `<https://example.com/font.woff2>;rel="preload";as="font";` +
			`crossorigin;media="screen";type="font/woff2"`
		if got := p.Header(); got != want {
			t.Errorf("p.Header() = %q, want %q", got, want)
		}
	})
}
