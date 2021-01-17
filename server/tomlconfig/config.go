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

// Package tomlconfig defines the TOML config for Web Packager HTTP Server.
package tomlconfig

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// TODO(banaag): Add the config for ACME.

// Config defines the TOML config.
// See cmd/webpkgserver/webpkgserver.example.toml for detail.
type Config struct {
	Listen    ListenConfig
	Server    ServerConfig
	SXG       SXGConfig
	Sign      SignConfig
	Processor ProcessorConfig
}

// ListenConfig represents the [Listen] section.
type ListenConfig struct {
	Host string
	Port int
	TLS  TLSConfig
}

// TLSConfig is part of ListenConfig.
type TLSConfig struct {
	PEMFile string
	KeyFile string
}

// ServerConfig represents the [Server] section.
type ServerConfig struct {
	DocPath      string `default:"/priv/doc"`
	CertPath     string `default:"/webpkg/cert"`
	ValidityPath string `default:"/webpkg/validity"`
	SignParam    string `default:"sign"`
}

// SXGConfig represents the [SXG] section.
type SXGConfig struct {
	Expiry             string `default:"168h"`
	JSExpiry           string `default:"24h"`
	CertURLBase        string `default:"/webpkg/cert"`
	ValidityURL        string `default:"/webpkg/validity"`
	KeepNonSXGPreloads bool
	Cert               SXGCertConfig
	ACME               SXGACMEConfig
}

// SXGCertConfig represents the [SXG.Cert] section.
type SXGCertConfig struct {
	PEMFile       string
	KeyFile       string
	CacheDir      string
	AllowTestCert bool
}

// SXGACMEConfig represents the [SXG.ACME] section.
type SXGACMEConfig struct {
	Enable            bool
	CSRFile           string
	DiscoveryURL      string
	Email             string
	EABKid            string
	EABHmac           string
	HTTPChallengePort int
	HTTPWebRootDir    string
	TLSChallengePort  int
	DNSProvider       string
}

// SignConfig represents the [[Sign]] sections.
type SignConfig []URLConfig

// URLConfig represents each of the [[Sign]] sections.
type URLConfig struct {
	Domain  string
	PathRE  string `default:".*"`
	QueryRE string `default:""`
}

// ProcessorConfig represents the [Processor] section.
type ProcessorConfig struct {
	SizeLimit  int `default:"4194304"`
	PreloadCSS bool
	PreloadJS  bool
}

// ReadFromFile reads a Config from filename. It also validates all fields
// and returns error if the validation fails.
func ReadFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfig(data)
}

// ParseConfig parses data into a Config. It also validates all fields and
// returns error if the validation fails.
func ParseConfig(data []byte) (*Config, error) {
	cfg := new(Config)
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if err := cfg.Verify(); err != nil {
		return nil, err
	}
	return cfg, nil
}
