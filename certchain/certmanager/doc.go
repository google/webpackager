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

/*
Package certmanager manages signed exchange certificates.

To get the certificate managed, the caller needs to specify the location of
the PEM file and calls the Manager's Start method at minimum:

	m := certmanager.NewManager(certmanager.Config{
		RawChainSource: certmanager.WatchCertFile(certmanager.WatchConfig{
			Path: "/path/to/your.pem",
		})
	})
	m.Start()
	defer m.Stop()

The above code lets the Manager check the PEM file every hour and retrieve
the OCSP response from the OCSP responder as needed.

The Manager caches the certificate chain and the OCSP response into a disk
if configured with a DiskCache. It enables them to be shared among multiple
processes, e.g. with other webpackager instances, in order to reduce the OCSP
responder's load.

Internals

Manager consists of two components: Producer and Cache. Producer produces
certificates and sends them to Manager continuously. Cache stores produced
certificates somewhere and possibly allows sharing them beyond the current
running process.

The certmanager package provides Augmentor as the canonical implementation
of Producer, which is composed of RawChainSource and OCSPRespSource. It looks
for the updated certificate chain and OCSP response at the right timings, and
uses NewAugmentedChain to turn them into an AugmentedChain. The "timings" are
controlled by RawChainSource and OCSPRespSource; see the GoDoc of those types
for details.
*/
// BUG(tomokinat): Manager writes to Cache, but certmanager does not reuse the
// cached information yet.
package certmanager
