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

type ProfileResource struct {
	Name string              `json:"name,omitempty" yaml:"name,omitempty"`
	Get  []ResourceOperation `json:"get,omitempty" yaml:"get,omitempty"`
	Set  []ResourceOperation `json:"set,omitempty" yaml:"set,omitempty"`
}

// String returns a JSON encoded string representation of the model
func (pr ProfileResource) String() string {
	out, err := json.Marshal(pr)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
