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
Package complexproc provides a factory of fully-featured processors, namely
NewComprehensiveProcessor. A ComprehensiveProcessor consists of prerequisite
checkers, preprocessors, main processors, and postprocessors.

The prerequisite checkers inspect each resource and verify it is eligible
for a signed exchange, provided by preverify.CheckPrerequisites.

The main processors include the key optimization logic and are defined on
a per Content-Type basis. The default ones just include an HTML processor
provided by the htmlproc package, applied to HTML (text/html) and XHTML
(application/xhtml+xml). The caller can change the behavior of the HTML
processor and/or add main processors for other media types through Config.
It is also possible to override or disable the default processors.

The preprocessors and the postprocessors include the logic applied to all
resources, before and after the main processors respectively. Some of them
are considered to be essential and always included. The caller can specify
additional preprocessors and postprocessors via Config; they are run after
the essential ones.
*/
package complexproc
