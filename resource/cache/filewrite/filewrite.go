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
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/webpackager/resource"
	"github.com/google/webpackager/resource/cache"
)

// NewFileWriteCache creates and initializes a new ResourceCache that also
// saves signed exchanges to files on the Store operations.
func NewFileWriteCache(config Config) cache.ResourceCache {
	return &fileWriteCache{config}
}

type fileWriteCache struct {
	Config
}

func (fsc *fileWriteCache) Lookup(req *http.Request) (*resource.Resource, error) {
	return fsc.BaseCache.Lookup(req)
}

func (fsc *fileWriteCache) Store(r *resource.Resource) error {
	if err := fsc.BaseCache.Store(r); err != nil {
		return err
	}
	if fsc.ExchangeMapping != nil && r.Exchange != nil {
		if err := write(fsc.ExchangeMapping, r, r.Exchange); err != nil {
			return err
		}
	}

	return nil
}

type writable interface {
	Write(w io.Writer) error
}

func write(mapping MappingRule, r *resource.Resource, data writable) error {
	path, err := mapping.Map(r)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return data.Write(file)
}
