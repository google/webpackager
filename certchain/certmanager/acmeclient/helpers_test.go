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

package acmeclient_test

import (
	"net"
	"testing"
)

func getFreeTCPPort(t *testing.T) int {
	t.Helper()

	s, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find a free tcp port: %v", err)
	}
	defer s.Close()

	return s.Addr().(*net.TCPAddr).Port
}
