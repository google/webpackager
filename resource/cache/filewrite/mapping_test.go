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
	"errors"
	"fmt"
	"testing"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/cache/filewrite"
)

func TestUseURLPathRules_Success(t *testing.T) {
	r := resource.NewResource(urlutil.MustParse("https://example.com/hello/"))
	r.PhysicalURL = urlutil.MustParse("https://example.com/hello/index.html")
	r.ValidityURL = urlutil.MustParse("https://example.com/hello/index.html.validity.1234567890")

	tests := []struct {
		name string
		rule filewrite.MappingRule
		want string
	}{
		{
			name: "UsePhysicalURLPath",
			rule: filewrite.UsePhysicalURLPath(),
			want: "hello/index.html",
		},
		{
			name: "UseValidityURLPath",
			rule: filewrite.UseValidityURLPath(),
			want: "hello/index.html.validity.1234567890",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.rule.Map(r)
			if err != nil {
				t.Fatalf("got error(%q), want success", err)
			}
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestUseURLPathRules_Error(t *testing.T) {
	tests := []struct {
		name string
		rule filewrite.MappingRule
		arg  *resource.Resource
	}{
		{
			name: "UsePhysicalURLPath_Unclean",
			rule: filewrite.UsePhysicalURLPath(),
			arg: &resource.Resource{
				PhysicalURL: urlutil.MustParse("https://example.com/hello/./index.html"),
			},
		},
		{
			name: "UsePhysicalURLPath_IsDir",
			rule: filewrite.UsePhysicalURLPath(),
			arg: &resource.Resource{
				PhysicalURL: urlutil.MustParse("https://example.com/hello/"),
			},
		},
		{
			name: "UseValidityURLPath_Unclean",
			rule: filewrite.UseValidityURLPath(),
			arg: &resource.Resource{
				ValidityURL: urlutil.MustParse("https://example.com/hello/./index.html.validity.1234567890"),
			},
		},
		{
			name: "UseValidityURLPath_IsDir",
			rule: filewrite.UseValidityURLPath(),
			arg: &resource.Resource{
				ValidityURL: urlutil.MustParse("https://example.com/hello/"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.rule.Map(test.arg)
			if err == nil {
				t.Fatal("got success, want non-nil error")
			}
		})
	}
}

func TestDecoratorRules(t *testing.T) {
	errDummy := errors.New("errDummy")

	tests := []struct {
		name string
		rule filewrite.MappingRule
		want string
		err  error
	}{
		{
			name: "AddBaseDir_Success",
			rule: filewrite.AddBaseDir(FixedMappingRule("hello/world.html"), "/tmp"),
			want: "/tmp/hello/world.html",
		},
		{
			name: "AddBaseDir_Error",
			rule: filewrite.AddBaseDir(ErrorMappingRule(errDummy), "/tmp"),
			err:  errDummy,
		},
		{
			name: "AddBaseDir_DevNull",
			rule: filewrite.AddBaseDir(filewrite.MapToDevNull(), "/tmp"),
			want: "",
		},
		{
			name: "AppendExt_Success",
			rule: filewrite.AppendExt(FixedMappingRule("hello/world.html"), ".sxg"),
			want: "hello/world.html.sxg",
		},
		{
			name: "AppendExt_Error",
			rule: filewrite.AppendExt(ErrorMappingRule(errDummy), ".sxg"),
			err:  errDummy,
		},
		{
			name: "AppendExt_DevNull",
			rule: filewrite.AppendExt(filewrite.MapToDevNull(), ".sxg"),
			want: "",
		},
		{
			name: "StripDir_Success",
			rule: filewrite.StripDir(FixedMappingRule("hello/world.html")),
			want: "world.html",
		},
		{
			name: "StripDir_Error",
			rule: filewrite.StripDir(ErrorMappingRule(errDummy)),
			err:  errDummy,
		},
		{
			name: "StripDir_DevNull",
			rule: filewrite.StripDir(filewrite.MapToDevNull()),
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.rule.Map(&resource.Resource{})
			if err != test.err {
				desc := func(err error) string {
					switch err {
					case nil:
						return "success"
					case errDummy:
						return "errDummy"
					default:
						return fmt.Sprintf("error(%q)", err)
					}
				}
				t.Errorf("got %s, want %s", desc(err), desc(test.err))
			}
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
