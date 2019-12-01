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

// Action describes state related to the capabilities of a device
type Action struct {
	Path      string     `json:"path,omitempty" yaml:"path,omitempty"`           // Path used by service for action on a device or sensor
	Responses []Response `json:"responses,omitempty" yaml:"responses,omitempty"` // Responses from get or put requests to service
	URL       string     `json:"url,omitempty" yaml:"url,omitempty"`             // Url for requests from command service
}

// String returns a JSON formatted string representation of the Action
func (a Action) String() string {
	out, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
