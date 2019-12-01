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
	"reflect"
	"strings"
)

// Response for a Get or Put request to a service
type Response struct {
	Code           string   `json:"code,omitempty" yaml:"code,omitempty"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	ExpectedValues []string `json:"expectedValues,omitempty" yaml:"expectedValues,omitempty"`
}

// String returns a JSON encoded string representation of the model
func (r Response) String() string {
	out, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func (r Response) Equals(r2 Response) bool {
	if strings.Compare(r.Code, r2.Code) != 0 {
		return false
	}
	if strings.Compare(r.Description, r2.Description) != 0 {
		return false
	}
	if len(r.ExpectedValues) != len(r2.ExpectedValues) {
		return false
	}
	if !reflect.DeepEqual(r.ExpectedValues, r2.ExpectedValues) {
		return false
	}
	return true

}
