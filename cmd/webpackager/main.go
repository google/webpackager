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

// webpackager is a command to "package" websites in accordance with
// https://github.com/WICG/webpackage/.
//
// See README.md for more information.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/webpackager"
	"github.com/google/webpackager/internal/multierror"
)

func run() error {
	flag.Parse()

	urls, err := getURLListFromFlags()
	if err != nil {
		return err
	}
	cfg, err := getConfigFromFlags()
	if err != nil {
		return err
	}
	pkg := webpackager.NewPackager(*cfg)

	for _, u := range urls {
		pkg.Run(u)
	}

	return pkg.Err()
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
}

func main() {
	if err := run(); err != nil {
		if errs, ok := err.(*multierror.MultiError); ok {
			for _, err := range errs.Errors {
				printError(err)
			}
		} else {
			printError(err)
		}

		os.Exit(1)
	}
}
