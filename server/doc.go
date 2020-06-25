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
Package server implements Web Packager HTTP Server (webpkgserver).

If you are interested in building and running webpkgserver as a binary,
see cmd/webpkgserver/README.md instead.

Basic Use

FromTOMLConfig creates a Server that can be used out of the box:

	c, err := tomlconfig.ReadFromFile("your.toml")
	if err != nil {
		log.Fatal(err)
	}
	s, err := server.FromTOMLConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	s.ListenAndServe() // Use ListenAndServeTLS to enable TLS.

Define Custom Parameters

If you want to define custom parameters/sections in TOML, define a struct
with a tomlconfig.Config embedded:

	type Config struct {
		tomlconfig.Config
		Foo FooConfig
	}

With the example above, your TOML config can contain the [Foo] section in
addition to the standard ones.

You need to call toml.Unmarshal by yourself. Also be sure to call Verify
on the tomlconfig.Config embedding; otherwise FromTOMLConfig may panic with
invalid config values.

	data, err := ioutil.ReadFile("your.toml")
	if err != nil {
		log.Fatal(err)
	}
	var c Config
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatal(err)
	}
	if err := c.Verify(); err != nil {
		log.Fatal(err)
	}
	s, err := server.FromTOMLConfig(&c.Config)
	if err != nil {
		log.Fatal(err)
	}
	// ... (mutate s.Packager and s.CertManager to apply FooConfig settings)
	s.ListenAndServe() // Use ListenAndServeTLS to enable TLS.

Handler Internals

Handler is composed of three child handlers: doc handler, cert handler, and
validity handler.

The doc handler produces a signed exchange for the given URL. The request
looks like:

	/priv/doc/https://example.com/index.html
	    -- or --
	/priv/doc?sign=https%3A%2F%2Fexample.com%2Findex.html

where "/priv/doc" and "sign" can be customized through DocPath and SignParam
in tomlconfig.ServerConfig respectively.

The cert handler serves AugmentedChains in the application/cert-chain+cbor
format. The request looks like:

	/webpkg/cert/47DEQpj8HBSa-_TImW+5JCeuQeRkm5NMpJWZG3hSuFUK

where "/webpkg/cert" can be customized through CertPath and "47DEQpj8..." is
an example of unique stable identifier, which is RawChain.Digest of the served
AugmentedChain.

The validity handler serves validity data. Currently, it constantly returns
an empty CBOR map (a single byte of 0xa0), which is interpreted as "no update
available." The request looks like:

    /webpkg/validity

where "/webpkg/validity" can be customized through ValidityPath. It does not
take any argument, such as the document URL, at this moment.
*/
package server
