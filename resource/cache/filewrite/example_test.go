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
	"fmt"
	"log"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/cache/filewrite"
)

func ExampleAddBaseDir() {
	urlParse := urlutil.MustParse // Like url.Parse, panicking on an error.

	r := resource.NewResource(urlParse("https://example.com/hello/world/"))
	r.PhysicalURL = urlParse("https://example.com/hello/world/index.html")
	r.ValidityURL = urlParse("https://example.com/hello/world/index.html.validity.1564617600")

	mapping := filewrite.AddBaseDir(filewrite.UsePhysicalURLPath(), "/tmp")

	got, err := mapping.Map(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	// Output:
	// /tmp/hello/world/index.html
}

func ExampleAppendExt() {
	urlParse := urlutil.MustParse // Like url.Parse, panicking on an error.

	r := resource.NewResource(urlParse("https://example.com/hello/world/"))
	r.PhysicalURL = urlParse("https://example.com/hello/world/index.html")
	r.ValidityURL = urlParse("https://example.com/hello/world/index.html.validity.1564617600")

	mapping := filewrite.AppendExt(filewrite.UsePhysicalURLPath(), ".sxg")

	got, err := mapping.Map(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	// Output:
	// hello/world/index.html.sxg
}

func ExampleStripDir() {
	urlParse := urlutil.MustParse // Like url.Parse, panicking on an error.

	r := resource.NewResource(urlParse("https://example.com/hello/world/"))
	r.PhysicalURL = urlParse("https://example.com/hello/world/index.html")
	r.ValidityURL = urlParse("https://example.com/hello/world/index.html.validity.1564617600")

	mapping := filewrite.StripDir(filewrite.UsePhysicalURLPath())

	got, err := mapping.Map(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	// Output:
	// index.html
}

func ExampleUsePhysicalURLPath() {
	urlParse := urlutil.MustParse // Like url.Parse, panicking on an error.

	r := resource.NewResource(urlParse("https://example.com/hello/world/"))
	r.PhysicalURL = urlParse("https://example.com/hello/world/index.html")
	r.ValidityURL = urlParse("https://example.com/hello/world/index.html.validity.1564617600")

	mapping := filewrite.UsePhysicalURLPath()

	got, err := mapping.Map(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	// Output:
	// hello/world/index.html
}

func ExampleUseValidityURLPath() {
	urlParse := urlutil.MustParse // Like url.Parse, panicking on an error.

	r := resource.NewResource(urlParse("https://example.com/hello/world/"))
	r.PhysicalURL = urlParse("https://example.com/hello/world/index.html")
	r.ValidityURL = urlParse("https://example.com/hello/world/index.html.validity.1564617600")

	mapping := filewrite.UseValidityURLPath()

	got, err := mapping.Map(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	// Output:
	// hello/world/index.html.validity.1564617600
}
