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

// Put models a put command in EdgeX
type Put struct {
	Action         `yaml:",inline"`
	ParameterNames []string `json:"parameterNames,omitempty" yaml:"parameterNames,omitempty"`
}

// String returns a JSON encoded string representation of the model
func (p Put) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// Append the associated value descriptors to the list
func (p *Put) AllAssociatedValueDescriptors(vdNames *map[string]string) {
	for _, pn := range p.ParameterNames {
		// Only add to the map if the value descriptor is NOT there
		if _, ok := (*vdNames)[pn]; !ok {
			(*vdNames)[pn] = pn
		}
	}
}
