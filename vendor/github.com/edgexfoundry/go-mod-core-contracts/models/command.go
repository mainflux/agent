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

// Command defines a specific read/write operation targeting a device
type Command struct {
	Timestamps  `yaml:",inline"`
	Id          string `json:"id" yaml:"id,omitempty"`     // Id is a unique identifier, such as a UUID
	Name        string `json:"name" yaml:"name,omitempty"` // Command name (unique on the profile)
	Get         Get    `json:"get" yaml:"get,omitempty"`   // Get Command
	Put         Put    `json:"put" yaml:"put,omitempty"`   // Put Command
	isValidated bool   // internal member used for validation check
}

// MarshalJSON implements the Marshaler interface. Empty strings will be null.
func (c Command) MarshalJSON() ([]byte, error) {
	test := struct {
		Timestamps
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"` // Command name (unique on the profile)
		Get  *Get   `json:"get,omitempty"`  // Get Command
		Put  *Put   `json:"put,omitempty"`  // Put Command
	}{
		Timestamps: c.Timestamps,
		Id:         c.Id,
		Name:       c.Name,
		Get:        &c.Get,
		Put:        &c.Put,
	}

	// Make empty structs nil pointers so they aren't marshaled
	if reflect.DeepEqual(c.Get, Get{}) {
		test.Get = nil
	}
	if reflect.DeepEqual(c.Put, Put{}) {
		test.Put = nil
	}

	return json.Marshal(test)
}

// UnmarshalJSON implements the Unmarshaler interface for the Command type
func (c *Command) UnmarshalJSON(data []byte) error {
	var err error
	a := new(struct {
		Timestamps `json:",inline"`
		Id         *string `json:"id"`
		Name       *string `json:"name"` // Command name (unique on the profile)
		Get        Get     `json:"get"`  // Get Command
		Put        Put     `json:"put"`  // Put Command
	})

	// Error with unmarshaling
	if err = json.Unmarshal(data, a); err != nil {
		return err
	}

	// Check nil fields
	if a.Id != nil {
		c.Id = *a.Id
	}
	if a.Name != nil {
		c.Name = *a.Name
	}
	c.Get = a.Get
	c.Put = a.Put
	c.Timestamps = a.Timestamps

	c.isValidated, err = c.Validate()

	return err
}

// Validate satisfies the Validator interface
func (c Command) Validate() (bool, error) {
	if !c.isValidated {
		if c.Name == "" {
			return false, NewErrContractInvalid("Name cannot be blank")
		}
		err := validate(c)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.isValidated, nil
}

/*
 * String() function for formatting
 */
func (c Command) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// AllAssociatedValueDescriptors will append all the associated value descriptors to the list
// associated by PUT command parameters and PUT/GET command return values
func (c *Command) AllAssociatedValueDescriptors(vdNames *map[string]string) {
	// Check and add Get value descriptors
	if &(c.Get) != nil {
		c.Get.AllAssociatedValueDescriptors(vdNames)
	}

	// Check and add Put value descriptors
	if &(c.Put) != nil {
		c.Put.AllAssociatedValueDescriptors(vdNames)
	}
}
