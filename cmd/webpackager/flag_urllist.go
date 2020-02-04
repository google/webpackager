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

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/google/webpackager/internal/customflag"
	multierror "github.com/hashicorp/go-multierror"
)

var (
	flagURL     = customflag.MultiString("url", `URL of an HTML page. Ignored when --url_file is given. (repeatable)`)
	flagURLFile = flag.String("url_file", "", `File to read the URL list from, or "-" to read it from stdin.`)
)

func getURLListFromFlags() ([]*url.URL, error) {
	unparsed, err := getURLStringList()
	if err != nil {
		return nil, err
	}
	if len(unparsed) == 0 {
		return nil, errors.New("no urls to process")
	}

	errs := new(multierror.Error)

	urls := make([]*url.URL, len(unparsed))
	for i, s := range unparsed {
		urls[i], err = url.Parse(s)
		if err != nil {
			if uerr, ok := err.(*url.Error); ok {
				err = uerr.Err
			}
			errs = multierror.Append(errs, fmt.Errorf("malformed url %q: %v", s, err))
		}
	}
	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}
	return urls, nil
}

func getURLStringList() ([]string, error) {
	if *flagURLFile == "" {
		return *flagURL, nil
	}
	if len(*flagURL) != 0 {
		return nil, errors.New("--url and --url_file may not be used together")
	}
	f, err := openFile(*flagURLFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLines(bufio.NewScanner(f))
}

func openFile(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return ioutil.NopCloser(os.Stdin), nil
	}
	return os.Open(filename)
}

func removeComment(s string) string {
	return strings.SplitN(s, "#", 2)[0]
}

func readLines(scanner *bufio.Scanner) ([]string, error) {
	lines := []string{}
	for scanner.Scan() {
		if s := strings.TrimSpace(removeComment(scanner.Text())); s != "" {
			lines = append(lines, s)
		}
	}
	return lines, scanner.Err()
}
