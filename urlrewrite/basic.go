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

package urlrewrite

import (
	"net/http"
	"net/url"
	"path"

	"github.com/google/webpackager/internal/urlutil"
)

// CleanPath normalizes the URL path. It basically applies path.Clean to
// eliminate multiple slashes, "." (current directory), and ".." (parent
// directory) in the path, but adds the trailing slash back.
func CleanPath() Rule {
	return &cleanPath{}
}

type cleanPath struct{}

func (r *cleanPath) Rewrite(u *url.URL, respHeader http.Header) {
	u.Path = urlutil.GetCleanPath(u)
}

// IndexRule appends indexFile to the path when it ends with a slash ("/")
// or otherwise likely represents a directory. indexFile is supposed to match
// the index file on the target server (typically "index.html").
func IndexRule(indexFile string) Rule {
	return &indexRule{indexFile}
}

type indexRule struct {
	indexFile string
}

func (r *indexRule) Rewrite(u *url.URL, respHeader http.Header) {
	if urlutil.IsDir(u) {
		u.Path = path.Join(u.Path, r.indexFile)
	}
}
