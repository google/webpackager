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

/*
Package acmeclient provides a RawChainSource to acquire a signed exchange
certificate using the ACME protocol. The ACME protocol allows a server to
obtain a certificate automatically, without any human intervention.
To learn about how it works, see https://letsencrypt.org/how-it-works/.

Client is the ACME client that behaves as a RawChainSource, typically used
with Manager or Augmentor. The consumer code is recommended to use it with
DiskCache, so the admin can retrieve the acquired certificate from the disk:

	m := certmanager.Manager(certmanager.Config{
		RawChainSource: acmeclient.NewClient(acmeclient.Config{
			...
		},
		OCSPRespSource: certmanager.DefaultOCSPClient,
		Cache: certmanager.NewDiskCache(certmanager.DiskCacheConfig{
			CertPath: "/tmp/acme/cert.pem",
			LockPath: "/tmp/acme/.lock",
		},
	})

Client is capable of responding to "challenges," which the ACME server poses
to the requestor to verify it has the control over the requested domain. See
Config for how to set up the challenge responder.
*/
package acmeclient
