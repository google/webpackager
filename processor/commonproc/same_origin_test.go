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
	"reflect"
	"testing"

	"github.com/google/webpackager/exchange/exchangetest"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/processor/commonproc"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/preload"
)

func plain(rawurl, as string) preload.Preload {
	r := resource.NewResource(urlutil.MustParse(rawurl))
	return preload.NewPlainPreload(r, as)
}

func TestApplySameOriginPolicy(t *testing.T) {
	resp := exchangetest.MakeEmptyResponse("https://example.org/")

	resp.Preloads = []preload.Preload{
		plain("https://example.org/assets/foo.css", preload.AsStyle),
		plain("https://example.com/assets/bar.css", preload.AsStyle),
		plain("https://example.org/assets/baz.css", preload.AsStyle),
		plain("https://example.org/assets/qux.js", preload.AsScript),
	}
	want := []preload.Preload{
		plain("https://example.org/assets/foo.css", preload.AsStyle),
		plain("https://example.org/assets/baz.css", preload.AsStyle),
		plain("https://example.org/assets/qux.js", preload.AsScript),
	}

	if err := commonproc.ApplySameOriginPolicy.Process(resp); err != nil {
		t.Fatalf("got error(%q), want success", err)
	}
	if !reflect.DeepEqual(resp.Preloads, want) {
		t.Errorf("resp.Preloads = %v, want %v", resp.Preloads, want)
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
