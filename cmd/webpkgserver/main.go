// Copyright 2020 Google LLC
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

// webpkgserver is a command to run Web Packager HTTP Server.
//
// See README.md for more information.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/webpackager/server"
	"github.com/google/webpackager/server/tomlconfig"
	multierror "github.com/hashicorp/go-multierror"
)

var (
	flagConfig = flag.String("config", "webpkgserver.toml", "Config TOML file.")
)

func run() error {
	flag.Parse()

	c, err := tomlconfig.ReadFromFile(*flagConfig)
	if err != nil {
		return err
	}
	s, err := server.FromTOMLConfig(c)
	if err != nil {
		return err
	}

	// Create a Listener by ourselves to show the precise listener address,
	// especially when Listen.Port is unspecified in TOML config.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	if s.TLSConfig == nil {
		log.Printf("Listening at %s", ln.Addr())
		return s.Serve(ln)
	} else {
		log.Printf("Listening TLS at %s", ln.Addr())
		return s.ServeTLS(ln, "", "")
	}
}

func printError(err error) {
	if me, ok := err.(*multierror.Error); ok {
		for _, err := range me.Errors {
			printError(err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	}
}

func main() {
	if err := run(); err != nil {
		printError(err)
		os.Exit(1)
	}
}
