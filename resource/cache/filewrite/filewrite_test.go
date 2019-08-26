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

package filewrite_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/cache"
	"github.com/google/webpackager/resource/cache/filewrite"
)

func TestStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fswriter_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	sxgBytes, err := ioutil.ReadFile("../../../testdata/sxg/standalone.sxg")
	if err != nil {
		t.Fatal(err)
	}
	sxg, err := signedexchange.ReadExchange(bytes.NewReader(sxgBytes))
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, sxg.RequestURI, nil)
	if err != nil {
		t.Fatal(err)
	}

	r := resource.NewResource(req.URL)
	r.SetExchange(sxg)

	t.Run("WriteToFile", func(t *testing.T) {
		tempFile := filepath.Join(tempDir, "standalone.sxg")
		cache := filewrite.NewFileWriteCache(filewrite.Config{
			BaseCache:       cache.NewOnMemoryCache(),
			ExchangeMapping: FixedMappingRule(tempFile)})

		if err := cache.Store(r); err != nil {
			t.Fatalf("cache.Store()  = error(%q), want success", err)
		}

		got, err := cache.Lookup(req)
		if err != nil {
			t.Fatalf("cache.Lookup() = error(%q), want success", err)
		}
		if got != r {
			t.Errorf("cache.Lookup() = <%p>, want <%p>", got, r)
		}

		gotBytes, err := ioutil.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("ioutil.ReadFile() = error(%q), want success", err)
		}
		if !bytes.Equal(gotBytes, sxgBytes) {
			t.Errorf("ioutil.ReadFile() = %q (%d bytes), want %q (%d bytes)",
				gotBytes, len(gotBytes), sxgBytes, len(sxgBytes))
		}
	})

	t.Run("WriteToNone", func(t *testing.T) {
		cache := filewrite.NewFileWriteCache(filewrite.Config{
			BaseCache:       cache.NewOnMemoryCache(),
			ExchangeMapping: filewrite.MapToDevNull()})

		if err := cache.Store(r); err != nil {
			t.Fatalf("cache.Store()  = error(%q), want success", err)
		}

		got, err := cache.Lookup(req)
		if err != nil {
			t.Fatalf("cache.Lookup() = error(%q), want success", err)
		}
		if got != r {
			t.Errorf("cache.Lookup() = <%p>, want <%p>", got, r)
		}
	})
}
