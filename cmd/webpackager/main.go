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
	multierror "github.com/hashicorp/go-multierror"
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
	date, err := getDateFromFlags()
	if err != nil {
		return err
	}

	pkg := webpackager.NewPackager(*cfg)
	errs := new(multierror.Error)

	for _, u := range urls {
		if err := pkg.Run(u, date); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs.ErrorOrNil()
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
