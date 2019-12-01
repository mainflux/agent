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
)

type ProfileProperty struct {
	Value PropertyValue `json:"value"`
	Units Units         `json:"units"`
}

// MarshalJSON implements the Marshaler interface
func (pp ProfileProperty) MarshalJSON() ([]byte, error) {
	test := struct {
		Value *PropertyValue `json:"value,omitempty"`
		Units *Units         `json:"units,omitempty"`
	}{
		Value: &pp.Value,
		Units: &pp.Units,
	}

	if reflect.DeepEqual(pp.Value, PropertyValue{}) {
		test.Value = nil
	}
	if reflect.DeepEqual(pp.Units, Units{}) {
		test.Units = nil
	}

	return json.Marshal(test)
}

// String returns a JSON encoded string representation of this ProfileProperty
func (pp ProfileProperty) String() string {
	out, err := json.Marshal(pp)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
