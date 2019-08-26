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
Package webpackager implements the control flow of Web Packager.

The code below illustrates the usage of this package:

	packager := webpackager.NewPackager(webpackager.Config{
		ExchangeFactory: &exchange.Factory{
			Version:      version.Version1b3,
			MIRecordSize: 4096,
			// ... (you need to set other fields)
		},
		ResourceCache: filewrite.NewFileWriteCache(&filewrite.Config{
			BaseCache: cache.NewOnMemoryCache(),
			ExchangeMapping: filewrite.AddBaseDir(
				filewrite.AppendExt(filewrite.UsePhysicalURLPath(), *flagSXGExt), *flagSXGDir),
		}),
	})
	for _, url := range urls {
		packager.Run(url)
	}
	if err := packager.Err(); err != nil {
		log.Fatal(err)
	}

Config allows you to change some behaviors of the Packager. packager.Run(url)
retrieves an HTTP response using FetchClient, processes it using Processor,
and turns it into a signed exchange using ExchangeFactory. Processor inspects
the HTTP response to see the eligibility for signed exchanges and manipulates
it to optimize the page loading. The generated signed exchanges are stored in
ResourceCache to prevent duplicates.

The code above sets just two parameters, ExchangeFactory and ResourceCache,
and uses the defaults for other parameters. With this setup, the packager
retrieves the content just through HTTP, applies the recommended set of
optimizations, generates signed exchanges compliant with the version b3, and
saves them in files named like "index.html.sxg" under "/tmp/sxg".

Config has a few more parameters. See its definition for the details.

You can also pass your own implementations to Config to inject custom logic
into the packaging flow. You could write, for example, a custom FetchClient
to retrieve the content from a database table instead of web servers, a custom
Processor or HTMLTask to apply your optimization techniques, a ResourceCache
to store the produced signed exchanges into another database table in addition
to a local drive, and so on.

The cmd/webpackager package provides a command line interface to execute the
packaging flow without writing the driver code on your own.
*/
package webpackager
