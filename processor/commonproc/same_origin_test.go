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

package commonproc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/processor/commonproc"
	"github.com/google/webpackager/resource/preload"
	"github.com/google/webpackager/resource/preload/preloadtest"
)

func TestApplySameOriginPolicy(t *testing.T) {
	resp := exchangetest.MakeEmptyResponse("https://example.org/")

	pl := preloadtest.NewPreloadForRawURL

	resp.Preloads = []*preload.Preload{
		pl("https://example.org/assets/foo.css", preload.AsStyle),
		pl("https://example.com/assets/bar.css", preload.AsStyle),
		pl("https://example.org/assets/baz.css", preload.AsStyle),
		pl("https://example.org/assets/qux.js", preload.AsScript),
	}
	want := []*preload.Preload{
		pl("https://example.org/assets/foo.css", preload.AsStyle),
		pl("https://example.org/assets/baz.css", preload.AsStyle),
		pl("https://example.org/assets/qux.js", preload.AsScript),
	}

	if err := commonproc.ApplySameOriginPolicy.Process(resp); err != nil {
		t.Fatalf("got error(%q), want success", err)
	}
	if diff := cmp.Diff(want, resp.Preloads); diff != "" {
		t.Errorf("resp.Preloads mismatch (-want +got):\n%s", diff)
	}
}

func TestApplySameOriginPolicy_Empty(t *testing.T) {
	resp := exchangetest.MakeEmptyResponse("https://example.org/")

	if err := commonproc.ApplySameOriginPolicy.Process(resp); err != nil {
		t.Fatalf("got error(%q), want success", err)
	}
	if len(resp.Preloads) != 0 {
		t.Errorf("resp.Preloads = %v, want []", resp.Preloads)
	}
}
