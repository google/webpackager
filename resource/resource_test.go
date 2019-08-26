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

package resource_test

import (
	"testing"

	"github.com/google/webpackager/exchange"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
)

func TestSetExchange(t *testing.T) {
	sxg, err := exchange.ReadExchangeFile("../testdata/sxg/standalone.sxg")
	if err != nil {
		panic(err)
	}
	integrity := "sha256-wvg1UzKYwDJYYcra5dOhakPaRGTAF5+saiIgmcro83E="

	r := resource.NewResource(urlutil.MustParse(sxg.RequestURI))
	if err := r.SetExchange(sxg); err != nil {
		t.Fatalf("r.SetExchange(sxg) = error(%q), want success", err)
	}
	if r.Exchange != sxg {
		t.Errorf("r.Exchange = <%p>, want <%p>", r.Exchange, sxg)
	}
	if r.Integrity != integrity {
		t.Errorf("r.Integrity = %q, want %q", r.Integrity, integrity)
	}
}

func TestAllowedAltSXGHeader(t *testing.T) {
	r := resource.NewResource(urlutil.MustParse("https://example.com/foo.js"))
	r.Integrity = "sha256-/DummyStringForSHA256ValueInBase64Encoding/="

	want := `<https://example.com/foo.js>;rel="allowed-alt-sxg";` +
		`header-integrity="sha256-/DummyStringForSHA256ValueInBase64Encoding/="`

	if got := r.AllowedAltSXGHeader(); got != want {
		t.Errorf("got %#q, want %#q", got, want)
	}
}
