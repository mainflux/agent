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

type Units struct {
	Type         string `json:"type,omitempty" yaml:"type,omitempty"`
	ReadWrite    string `json:"readWrite,omitempty" yaml:"readWrite,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
}

// String returns a JSON encoded string representation of the model
func (u Units) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
