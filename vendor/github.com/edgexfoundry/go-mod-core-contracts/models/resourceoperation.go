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

type ResourceOperation struct {
	Index          string            `json:"index" yaml:"index,omitempty"`
	Operation      string            `json:"operation" yaml:"operation,omitempty"`
	Object         string            `json:"object" yaml:"object,omitempty"`                 // Deprecated
	DeviceResource string            `json:"deviceResource" yaml:"deviceResource,omitempty"` // The replacement of Object field
	Parameter      string            `json:"parameter" yaml:"parameter,omitempty"`
	Resource       string            `json:"resource" yaml:"resource,omitempty"`           // Deprecated
	DeviceCommand  string            `json:"deviceCommand" yaml:"deviceCommand,omitempty"` // The replacement of Resource field
	Secondary      []string          `json:"secondary" yaml:"secondary,omitempty"`
	Mappings       map[string]string `json:"mappings" yaml:"mappings,omitempty"`
	isValidated    bool              // internal member used for validation check
}

// MarshalJSON returns a JSON encoded byte representation of the model and performs custom autofill
func (ro ResourceOperation) MarshalJSON() ([]byte, error) {
	test := struct {
		Index          string             `json:"index,omitempty"`
		Operation      string             `json:"operation,omitempty"`
		Object         string             `json:"object,omitempty"`
		DeviceResource string             `json:"deviceResource,omitempty"`
		Parameter      string             `json:"parameter,omitempty"`
		Resource       string             `json:"resource,omitempty"`
		DeviceCommand  string             `json:"deviceCommand,omitempty"`
		Secondary      []string           `json:"secondary,omitempty"`
		Mappings       *map[string]string `json:"mappings,omitempty"`
	}{
		Index:          ro.Index,
		Operation:      ro.Operation,
		Object:         ro.Object,
		DeviceResource: ro.DeviceResource,
		Parameter:      ro.Parameter,
		Resource:       ro.Resource,
		DeviceCommand:  ro.DeviceCommand,
		Secondary:      ro.Secondary,
		Mappings:       &ro.Mappings,
	}

	// Empty maps are nil
	if len(ro.Mappings) == 0 {
		test.Mappings = nil
	}

	if ro.DeviceResource != "" {
		test.Object = ro.DeviceResource
	} else if ro.Object != "" {
		test.Object = ro.Object
		test.DeviceResource = ro.Object
	}

	if ro.DeviceCommand != "" {
		test.Resource = ro.DeviceCommand
	} else if ro.Resource != "" {
		test.DeviceCommand = ro.Resource
	}

	return json.Marshal(test)
}

// UnmarshalJSON implements the Unmarshaler interface for the ResourceOperation type
func (ro *ResourceOperation) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Index          *string           `json:"index"`
		Operation      *string           `json:"operation"`
		Object         *string           `json:"object"`
		DeviceResource *string           `json:"deviceResource"`
		Parameter      *string           `json:"parameter"`
		Resource       *string           `json:"resource"`
		DeviceCommand  *string           `json:"deviceCommand"`
		Secondary      []string          `json:"secondary"`
		Mappings       map[string]string `json:"mappings"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Check nil fields
	if a.Index != nil {
		ro.Index = *a.Index
	}
	if a.Operation != nil {
		ro.Operation = *a.Operation
	}
	if a.DeviceResource != nil {
		ro.DeviceResource = *a.DeviceResource
		ro.Object = *a.DeviceResource
	} else if a.Object != nil {
		ro.Object = *a.Object
		ro.DeviceResource = *a.Object
	}
	if a.Parameter != nil {
		ro.Parameter = *a.Parameter
	}
	if a.DeviceCommand != nil {
		ro.DeviceCommand = *a.DeviceCommand
		ro.Resource = *a.DeviceCommand
	} else if a.Resource != nil {
		ro.Resource = *a.Resource
		ro.DeviceCommand = *a.Resource
	}
	ro.Secondary = a.Secondary
	ro.Mappings = a.Mappings

	ro.isValidated, err = ro.Validate()

	return err
}

// Validate satisfies the Validator interface
func (ro ResourceOperation) Validate() (bool, error) {
	if !ro.isValidated {
		if ro.Object == "" && ro.DeviceResource == "" {
			return false, NewErrContractInvalid("Object and DeviceResource are both blank")
		}
		err := validate(ro)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return ro.isValidated, nil
}

// String returns a JSON encoded string representation of the model
func (ro ResourceOperation) String() string {
	out, err := json.Marshal(ro)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
