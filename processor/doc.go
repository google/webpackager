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
Package processor defines the Processor interface.

Processors manipulate an exchange.Response to optimize the page loading,
especially with privacy-preserving prefetch. The canonical example is to
detect subresources from the payload (usually an HTML document) and add the
Link headers so that they are preloaded together. Processors typically just
manipulate the headers, but aggressive ones may also modify the payload.

Processors can also serve as checkers. Such Processors do not mutate the
exchange.Response in any way. Instead they inspect the exchange.Response
and report an error if the response does not meet some criteria, e.g. when
the response has HTTP status code other than 200 (OK). The preverify package
provides processors falling into this category.

This package also provides interfaces to combine multiple Processors and use
them as a single Processor: see MultiplexedProcessor and SequentialProcessor
for details.

The defaultproc package provides a processor that can be used out of the box.
*/
package processor
