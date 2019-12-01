/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package models

import "encoding/json"

// DescribedObject is a hold-over from the Java conversion and is supposed to represent inheritance whereby a type
// with a Description property IS A DescribedObject. However since there is no inheritance in Go, this should be
// eliminated and the Description property moved to the relevant types. 4 types currently use this.
type DescribedObject struct {
	Timestamps  `yaml:",inline"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"` // Description. Capicé?
}

// String returns a JSON formatted string representation of this DescribedObject
func (o DescribedObject) String() string {
	out, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
