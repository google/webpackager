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

package fetchtest_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/webpackager/fetch/fetchtest"
)

func TestFetchClient(t *testing.T) {
	testText := []byte("Hello, world!")

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(
		"dummy.test/hello.txt",
		func(w http.ResponseWriter, req *http.Request) {
			w.Write(testText)
		},
	)

	server := httptest.NewTLSServer(serveMux)
	defer server.Close()

	client := fetchtest.NewFetchClient(server)

	req, err := http.NewRequest(http.MethodGet, "https://dummy.test/hello.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("got error(%q), want success", err)
	}

	t.Run("Requests", func(t *testing.T) {
		want := []*http.Request{req}
		if got := client.Requests(); !reflect.DeepEqual(got, want) {
			t.Fatalf("client.Requests() = %v, want %v", got, want)
		}
	})
	t.Run("Response", func(t *testing.T) {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("r.StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
		}
		got, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("resp.Body = error(%q), want success", err)
		}
		if !bytes.Equal(got, testText) {
			t.Errorf("resp.Body = %q, want %q", got, testText)
		}
	})
}
