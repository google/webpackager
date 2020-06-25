// Copyright 2019 Google LLC
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

package certmanager_test

import (
	"syscall"
	"time"

	"github.com/google/webpackager/certchain/certmanager"
)

// The following code creates a new LocalCertFile that reads cert.pem
// (roughly) every other hours.
func ExampleLocalCertFile_hourly() {
	c := certmanager.LocalCertFileConfig{
		Path:        "cert.pem",
		FetchTiming: certmanager.FetchAtIntervals(2 * time.Hour),
	}
	_ = certmanager.NewLocalCertFile(c)
}

// The following code creates a new LocalCertFile that reads cert.pem
// every time the running process receives the USR1 signal.
func ExampleLocalCertFile_signal() {
	c := certmanager.LocalCertFileConfig{
		Path:        "cert.pem",
		FetchTiming: certmanager.FetchOnSignal(syscall.SIGUSR1),
	}
	_ = certmanager.NewLocalCertFile(c)
}
