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

/*
Package urlrewrite reproduces server-side URL rewrite logic.

Web servers sometimes "rewrite" request URLs internally, without issuing
redirects to clients. For example, many web servers append the index file
(e.g. "index.html" and "default.asp") to URL paths pointing to a directory.
"https://www.example.com/" thus may be rewritten, say, to
"https://www.example.com/index.html". Another usage is content negotiation:
some web sites are configured, for instance, to rewrite "*.jpg" to "*.webp"
and serve WebP images to supported clients.

Package urlrewrite implements some simple rewrite rules, aiming to provide
reasonable approximates for static web servers (which simply serve static
files under the document root). It also defines interface that allows you to
(re)implement custom rules and combine them with the existing ones.

The rewritten URLs are used to determine the file to write signed exchanges
to in the filesystem, determine the validity URL (where the validity data is
served), and so on. It especially helps support multiple resources served at
a single URL, e.g. with content negotiation. It also helps the URL path
better reflect the physical location of the resource on the server, which in
turn helps signed exchanges and validity data produced in the same directory
as the original resource.

References

https://www.igvita.com/2013/05/01/deploying-webp-via-accept-content-negotiation/
(Deploying WebP via Accept Content Negotiation)
*/
package urlrewrite
