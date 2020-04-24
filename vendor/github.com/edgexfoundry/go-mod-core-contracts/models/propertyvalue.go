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

const (
	// Base64Encoding : the float value is represented in Base64 encoding
	Base64Encoding = "Base64"
	// ENotation : the float value is represented in eNotation
	ENotation = "eNotation"
)

type PropertyValue struct {
	Type          string `json:"type,omitempty" yaml:"type,omitempty"`                 // ValueDescriptor Type of property after transformations
	ReadWrite     string `json:"readWrite,omitempty" yaml:"readWrite,omitempty"`       // Read/Write Permissions set for this property
	Minimum       string `json:"minimum,omitempty" yaml:"minimum,omitempty"`           // Minimum value that can be get/set from this property
	Maximum       string `json:"maximum,omitempty" yaml:"maximum,omitempty"`           // Maximum value that can be get/set from this property
	DefaultValue  string `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"` // Default value set to this property if no argument is passed
	Size          string `json:"size,omitempty" yaml:"size,omitempty"`                 // Size of this property in its type  (i.e. bytes for numeric types, characters for string types)
	Mask          string `json:"mask,omitempty" yaml:"mask,omitempty"`                 // Mask to be applied prior to get/set of property
	Shift         string `json:"shift,omitempty" yaml:"shift,omitempty"`               // Shift to be applied after masking, prior to get/set of property
	Scale         string `json:"scale,omitempty" yaml:"scale,omitempty"`               // Multiplicative factor to be applied after shifting, prior to get/set of property
	Offset        string `json:"offset,omitempty" yaml:"offset,omitempty"`             // Additive factor to be applied after multiplying, prior to get/set of property
	Base          string `json:"base,omitempty" yaml:"base,omitempty"`                 // Base for property to be applied to, leave 0 for no power operation (i.e. base ^ property: 2 ^ 10)
	Assertion     string `json:"assertion,omitempty" yaml:"assertion,omitempty"`       // Required value of the property, set for checking error state.  Failing an assertion condition will mark the device with an error state
	Precision     string `json:"precision,omitempty" yaml:"precision,omitempty"`
	FloatEncoding string `json:"floatEncoding,omitempty" yaml:"floatEncoding,omitempty"` // FloatEncoding indicates the representation of floating value of reading.  It should be 'Base64' or 'eNotation'
	MediaType     string `json:"mediaType,omitempty" yaml:"mediaType,omitempty"`
}

// String returns a JSON encoded string representation of the model
func (pv PropertyValue) String() string {
	out, err := json.Marshal(pv)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
