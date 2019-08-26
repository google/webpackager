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

package filewrite

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/resource"
)

var (
	errBadPhysicalURL = errors.New(
		"filewrite: PhysicalURL is unclean or missing a filename")
	errBadValidityURL = errors.New(
		"filewrite: ValidityURL is unclean or missing a filename")
)

// MappingRule defines the rule of mapping Resources into files.
// More precisely, MappingRule determines the file to store the signed
// exchange for each Resource.
//
// In future, MappingRule will also be used for the validity data.
type MappingRule interface {
	// Map returns the path to the target file. The returned path can be
	// empty, in which case the data is not written to any file.
	Map(r *resource.Resource) (string, error)
}

// MapToDevNull returns a MappingRule to write no files.
func MapToDevNull() MappingRule {
	return &mapToDevNull{}
}

type mapToDevNull struct{}

func (*mapToDevNull) Map(r *resource.Resource) (string, error) {
	return "", nil
}

// UsePhysicalURLPath returns a MappingRule to use PhysicalURL.Path.
// The leading slash ("/") is stripped so the returned path becomes relative
// from the current working directory, not the root directory.
//
// PhysicalURL.Path must be cleaned (e.g. no "." or ".." elements) and have
// a filename. The UsePhysicalURLPath mapping returns an error otherwise.
func UsePhysicalURLPath() MappingRule {
	return &usePhysicalURLPath{}
}

type usePhysicalURLPath struct{}

func (*usePhysicalURLPath) Map(r *resource.Resource) (string, error) {
	u := r.PhysicalURL
	if urlutil.IsDir(u) || u.Path != urlutil.GetCleanPath(u) {
		return "", errBadPhysicalURL
	}
	return filepath.FromSlash(u.Path[1:]), nil
}

// UseValidityURLPath returns a MappingRule to use ValidityURL.Path.
// The leading slash ("/") is stripped so the returned path becomes relative
// from the current working directory, not the root directory.
//
// ValidityURL.Path must be cleaned (e.g. no "." or ".." elements) and have
// a filename. The UseValidityURLPath mapping returns an error otherwise.
func UseValidityURLPath() MappingRule {
	return &useValidityURLPath{}
}

type useValidityURLPath struct{}

func (*useValidityURLPath) Map(r *resource.Resource) (string, error) {
	u := r.ValidityURL
	if urlutil.IsDir(u) || u.Path != urlutil.GetCleanPath(u) {
		return "", errBadValidityURL
	}
	return filepath.FromSlash(u.Path[1:]), nil
}

// AddBaseDir returns a new MappingRule that calls rule.Map then prepends
// dir to the returned path.
func AddBaseDir(rule MappingRule, dir string) MappingRule {
	if dir == "" {
		return rule
	}
	return &addBaseDir{rule, dir}
}

type addBaseDir struct {
	base MappingRule
	dir  string
}

func (rule *addBaseDir) Map(r *resource.Resource) (string, error) {
	path, err := rule.base.Map(r)
	if path == "" || err != nil {
		return "", err
	}
	return filepath.Join(rule.dir, path), nil
}

// AppendExt returns a new MappingRule that calls rule.Map then appends
// ext to the returned path. ext usually starts with a period (e.g. ".sxg")
// and is known as a file extension of a file suffix.
//
// AppendExt panics if ext contains a directory separator ("/").
func AppendExt(rule MappingRule, ext string) MappingRule {
	if ext == "" {
		return rule
	}
	if strings.IndexRune(ext, filepath.Separator) >= 0 {
		panic("invalid extension")
	}
	return &appendExt{rule, ext}
}

type appendExt struct {
	base MappingRule
	ext  string
}

func (rule *appendExt) Map(r *resource.Resource) (string, error) {
	path, err := rule.base.Map(r)
	if path == "" || err != nil {
		return "", err
	}
	return (path + rule.ext), nil
}

// StripDir returns a new MappingRule that calls rule.Map then eliminates the
// directory part (anything but the last element) from the returned path.
func StripDir(rule MappingRule) MappingRule {
	return &stripDir{rule}
}

type stripDir struct {
	base MappingRule
}

func (rule *stripDir) Map(r *resource.Resource) (string, error) {
	path, err := rule.base.Map(r)
	if path == "" || err != nil {
		return "", err
	}
	return filepath.Base(path), nil
}
