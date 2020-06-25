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

package tomlconfig

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

var (
	errEmpty = errors.New("must be set non-empty")
	errRange = errors.New("value out of range")
)

// Verify validates all fields in c. ReadFromFile and ParseConfig calls Verify
// internally, so you do not usually have to call it by yourself.
func (c *Config) Verify() error {
	var errs *multierror.Error

	if err := c.Listen.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("Listen", err))
	}
	if err := c.Server.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("Server", err))
	}
	if err := c.SXG.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("SXG", err))
	}
	if err := c.Sign.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("Sign", err))
	}
	if err := c.Processor.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("Processor", err))
	}

	return errs.ErrorOrNil() // TODO(yuizumi): Format it better.
}

func (c *ListenConfig) verify() error {
	var errs *multierror.Error

	if c.Port < 0 || c.Port >= 65536 {
		errs = multierror.Append(errs, wrapError("Port", errRange))
	}
	if err := c.TLS.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("TLS", err))
	}

	return errs.ErrorOrNil()
}

func (c *TLSConfig) verify() error {
	if c.PEMFile != "" && c.KeyFile == "" {
		return errors.New("PEMFile specified without KeyFile")
	}
	if c.KeyFile != "" && c.PEMFile == "" {
		return errors.New("KeyFile specified without PEMFile")
	}
	return nil
}

func (c *ServerConfig) verify() error {
	var errs *multierror.Error

	if err := verifyServePath(c.DocPath); err != nil {
		errs = multierror.Append(errs, wrapError("DocPath", err))
	}
	if err := verifyServePath(c.CertPath); err != nil {
		errs = multierror.Append(errs, wrapError("CertPath", err))
	}
	if err := verifyServePath(c.ValidityPath); err != nil {
		errs = multierror.Append(errs, wrapError("ValidityPath", err))
	}
	if err := verifyParamName(c.SignParam); err != nil {
		errs = multierror.Append(errs, wrapError("SignParam", err))
	}

	return errs.ErrorOrNil()
}

func (c *SXGConfig) verify() error {
	var errs *multierror.Error

	if _, err := parseExpiry(c.Expiry); err != nil {
		errs = multierror.Append(errs, wrapError("Expiry", err))
	}
	if _, err := parseJSExpiry(c.JSExpiry); err != nil {
		errs = multierror.Append(errs, wrapError("JSExpiry", err))
	}
	if err := verifyURL(c.CertURLBase); err != nil {
		errs = multierror.Append(errs, wrapError("CertURLBase", err))
	}
	if err := verifyURL(c.ValidityURL); err != nil {
		errs = multierror.Append(errs, wrapError("ValidityURL", err))
	}

	if err := c.Cert.verify(); err != nil {
		errs = multierror.Append(errs, wrapError("Cert", err))
	}

	return errs.ErrorOrNil()
}

func (c *SXGCertConfig) verify() error {
	var errs *multierror.Error

	if c.PEMFile == "" {
		errs = multierror.Append(errs, wrapError("PEMFile", errEmpty))
	}
	if c.KeyFile == "" {
		errs = multierror.Append(errs, wrapError("KeyFile", errEmpty))
	}

	return errs.ErrorOrNil()
}

func (c SignConfig) verify() error {
	var errs *multierror.Error

	for i, uc := range c {
		if err := uc.verify(); err != nil {
			errs = multierror.Append(errs, wrapError(fmt.Sprintf("[%d]", i), err))
		}
	}

	return errs.ErrorOrNil()
}

func (c *URLConfig) verify() error {
	var errs *multierror.Error

	// TODO(yuizumi): Restrict to ASCII?
	if c.Domain == "" {
		errs = multierror.Append(errs, wrapError("Domain", errEmpty))
	}

	if _, err := regexp.Compile(c.PathRE); err != nil {
		errs = multierror.Append(errs, wrapError("PathRE", err))
	}
	if _, err := regexp.Compile(c.QueryRE); err != nil {
		errs = multierror.Append(errs, wrapError("QueryRE", err))
	}

	return errs.ErrorOrNil()
}

func (c *ProcessorConfig) verify() error {
	var errs *multierror.Error

	if c.SizeLimit <= 0 {
		errs = multierror.Append(errs, wrapError("SizeLimit", errRange))
	}

	return errs.ErrorOrNil()
}

func verifyParamName(value string) error {
	if value == "" {
		return errEmpty
	}
	return nil // TODO(yuizumi): Restrict to alphanumerics?
}

func verifyServePath(value string) error {
	if value == "" {
		return errEmpty
	}
	// Now we check the path is normalized using path.Clean, but we want
	// to allow the trailing slash while path.Clean removes it. To align
	// them, remove the slash here except for special cases.
	if value != "/" && value != "//" {
		value = strings.TrimSuffix(value, "/")
	}

	if value != path.Clean(value) {
		return errors.New("must be absolute and normalized")
	}
	if !strings.HasPrefix(value, "/") {
		return errors.New("must be absolute and normalized")
	}

	return nil
}

func verifyURL(value string) error {
	if value == "" {
		return errEmpty
	}

	u, err := url.Parse(value)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "":
		if u.Host != "" {
			return errors.New("must be https:// or absolute path")
		}
		if !strings.HasPrefix(u.Path, "/") {
			return errors.New("must be https:// or absolute path")
		}
		return verifyServePath(u.Path)
	case "https":
		return verifyServePath(u.Path)
	default:
		return errors.New("must be https:// or absolute path")
	}
}
