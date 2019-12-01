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

type Get struct {
	Action `json:",omitempty" yaml:",inline"`
}

/*
 * To String function for Get Struct
 */
func (g Get) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// Append the associated value descriptors to the list
func (g *Get) AllAssociatedValueDescriptors(vdNames *map[string]string) {
	for _, r := range g.Action.Responses {
		for _, ev := range r.ExpectedValues {
			// Only add to the map if the value is not there
			if _, ok := (*vdNames)[ev]; !ok {
				(*vdNames)[ev] = ev
			}
		}
	}
}
