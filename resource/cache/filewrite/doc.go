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
Package filewrite provides ResourceCache that also saves signed exchanges
to files on the Store operations to the cache. The ResourceCache works as
a wrapper around another ResourceCache and uses a MappingRule to locate the
files to write the signed exchanges to. Here is an example:

	exampleCache := filewrite.NewFileWriteCache(filewrite.Config{
		BaseCache: cache.NewOnMemoryCache(),
		ExchangeMapping: filewrite.AddBaseDir(
			filewrite.AppendExt(filewrite.UsePhysicalURLPath(), ".sxg"),
			"/tmp/sxg",
		),
	})

exampleCache uses an on-memory cache for the underlying cache. It saves
signed exchanges under /tmp/sxg, to the location parallel to the URL path,
and with the .sxg extension. For example, the signed exchange for:

	https://www.example.com/hello/world/index.html

would be saved to:

	/tmp/sxg/hello/world/index.html.sxg
*/
package filewrite
