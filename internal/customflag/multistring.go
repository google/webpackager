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

package customflag

import (
	"flag"
	"fmt"
)

type multiStringValue []string

func (msv *multiStringValue) String() string {
	if len(*msv) == 0 {
		return ""
	}
	return fmt.Sprint(*msv)
}

func (msv *multiStringValue) Set(value string) error {
	*msv = append(*msv, value)
	return nil
}

func (msv *multiStringValue) Get() interface{} {
	return []string(*msv)
}

// NewMultiStringValue returns a flag.Value for a []string flag. The argument p
// points to a []string variable in which to store the value of the flag.
//
// The default value is a nil slice. Each occurrence of the flag in the command
// line appends the value to the slice.
func NewMultiStringValue(p *[]string) flag.Value {
	*p = nil
	return (*multiStringValue)(p)
}

// MultiString defines a []string flag with specified name and usage string.
// The returned value is the address of a []string variable that stores the
// value of the flag.
//
// The default value is a nil slice. Each occurrence of the flag in the command
// line appends the value to the slice.
func MultiString(name string, usage string) *[]string {
	p := new([]string)
	flag.Var(NewMultiStringValue(p), name, usage)
	return p
}

// MultiStringVar defines a []string flag with specified name and usage string.
// The argument p points to a []string variable in which to store the value of
// the flag.
//
// The default value is a nil slice. Each occurrence of the flag in the command
// line appends the value to the slice.
func MultiStringVar(p *[]string, name string, usage string) {
	flag.Var(NewMultiStringValue(p), name, usage)
}
