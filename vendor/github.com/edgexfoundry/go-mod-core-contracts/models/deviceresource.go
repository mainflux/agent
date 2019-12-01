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

// DeviceResource represents a value on a device that can be read or written
type DeviceResource struct {
	Description string            `json:"description" yaml:"description,omitempty"`
	Name        string            `json:"name" yaml:"name,omitempty"`
	Tag         string            `json:"tag" yaml:"tag,omitempty"`
	Properties  ProfileProperty   `json:"properties" yaml:"properties"`
	Attributes  map[string]string `json:"attributes" yaml:"attributes,omitempty"`
}

// MarshalJSON implements the Marshaler interface in order to make empty strings null
func (do DeviceResource) MarshalJSON() ([]byte, error) {
	test := struct {
		Description string             `json:"description,omitempty"`
		Name        string             `json:"name,omitempty"`
		Tag         string             `json:"tag,omitempty"`
		Properties  *ProfileProperty   `json:"properties,omitempty"`
		Attributes  *map[string]string `json:"attributes,omitempty"`
	}{
		Description: do.Description,
		Name:        do.Name,
		Tag:         do.Tag,
		Properties:  &do.Properties,
	}

	// Empty maps are null
	if len(do.Attributes) > 0 {
		test.Attributes = &do.Attributes
	}
	if reflect.DeepEqual(do.Properties, ProfileProperty{}) {
		test.Properties = nil
	}

	return json.Marshal(test)
}

/*
 * To String function for DeviceResource
 */
func (do DeviceResource) String() string {
	out, err := json.Marshal(do)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
