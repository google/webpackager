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

package vprule

import (
	"log"
	"mime"
	"time"

	"github.com/google/webpackager/exchange"
)

// PerContentType specifies a Rule per media type. rules is a map from
// media types to Rules; ruleElse is the Rule applied to other media types.
// The map keys should be all in lowercase and include no media parameters
// (e.g. "text/html", not "text/HTML" or "text/html; charset=utf-8").
//
// PerContentType looks for a rule applicable to the resp's Content-Type
// first. If there is none, PerContentType also looks for a rule for each
// Webpackager-Sub-Content-Type (resp.ExtraData[exchange.SubContentType]).
// If there is still no rule to apply, PerContentType applies ruleElse.
func PerContentType(rules map[string]Rule, ruleElse Rule) Rule {
	return &perContentType{rules, ruleElse}
}

type perContentType struct {
	rules    map[string]Rule
	ruleElse Rule
}

func (p *perContentType) Get(resp *exchange.Response, date time.Time) exchange.ValidPeriod {
	if r := p.lookup(resp.Header.Get("Content-Type")); r != nil {
		return r.Get(resp, date)
	}
	for _, sct := range resp.ExtraData[exchange.SubContentType] {
		if r := p.lookup(sct); r != nil {
			return r.Get(resp, date)
		}
	}
	return p.ruleElse.Get(resp, date)
}

func (p *perContentType) lookup(mimeType string) Rule {
	if mimeType == "" {
		return nil
	}
	mediaType, _, err := mime.ParseMediaType(mimeType)
	if err != nil && err != mime.ErrInvalidMediaParameter {
		log.Printf("warning: invalid MIME type %q: %v", mimeType, err)
		return nil
	}
	return p.rules[mediaType]
}
