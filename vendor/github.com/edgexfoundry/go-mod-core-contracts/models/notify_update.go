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
	"fmt"
)

type NotifyUpdate struct {
	Name        string `json:"name,omitempty"`
	Operation   string `json:"operation,omitempty"`
	isValidated bool   // internal member used for validation check
}

// UnmarshalJSON implements the Unmarshaler interface for the NotifyUpdate type
func (n *NotifyUpdate) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Name      *string `json:"name"`
		Operation *string `json:"operation"`
	}
	a := Alias{}

	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	//Nillable fields
	if a.Name != nil {
		n.Name = *a.Name
	}
	if a.Operation != nil {
		n.Operation = *a.Operation
	}

	n.isValidated, err = n.Validate()

	return err
}

// Validate satisfies the Validator interface
func (n NotifyUpdate) Validate() (bool, error) {
	if !n.isValidated {
		if n.Name == "" || n.Operation == "" {
			return false, NewErrContractInvalid("Name and Operation must both have a value")
		}
		if n.Operation != NotifyUpdateAdd &&
			n.Operation != NotifyUpdateUpdate &&
			n.Operation != NotifyUpdateDelete {
			return false, NewErrContractInvalid(fmt.Sprintf("Invalid value for operation %s", n.Operation))
		}
	}
	return true, nil
}
