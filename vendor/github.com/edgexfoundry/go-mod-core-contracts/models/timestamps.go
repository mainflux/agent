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

import (
	"encoding/json"
)

type Timestamps struct {
	Created  int64 `json:"created,omitempty" yaml:"created,omitempty"`
	Modified int64 `json:"modified,omitempty" yaml:"modified,omitempty"`
	Origin   int64 `json:"origin,omitempty" yaml:"origin,omitempty"`
}

// String returns a JSON encoded string representation of the model
func (ts *Timestamps) String() string {
	out, err := json.Marshal(ts)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

/*
 * Compare the Created of two objects to determine given is newer
 */
func (ts *Timestamps) compareTo(i Timestamps) int {
	if i.Created > ts.Created {
		return 1
	}
	return -1
}
