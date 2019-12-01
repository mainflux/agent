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
	"regexp"
)

// defaultValueDescriptorFormat defines the default formatting value used with creating a ValueDescriptor from a DeviceResource.
const defaultValueDescriptorFormat = "%s"

/*
 * Value Descriptor Struct
 */
type ValueDescriptor struct {
	Id            string      `json:"id,omitempty"`
	Created       int64       `json:"created,omitempty"`
	Description   string      `json:"description,omitempty"`
	Modified      int64       `json:"modified,omitempty"`
	Origin        int64       `json:"origin,omitempty"`
	Name          string      `json:"name,omitempty"`
	Min           interface{} `json:"min,omitempty"`
	Max           interface{} `json:"max,omitempty"`
	DefaultValue  interface{} `json:"defaultValue,omitempty"`
	Type          string      `json:"type,omitempty"`
	UomLabel      string      `json:"uomLabel,omitempty"`
	Formatting    string      `json:"formatting,omitempty"`
	Labels        []string    `json:"labels,omitempty"`
	MediaType     string      `json:"mediaType,omitempty"`
	FloatEncoding string      `json:"floatEncoding,omitempty"`
	isValidated   bool        // internal member used for validation check
}

// UnmarshalJSON implements the Unmarshaler interface for the ValueDescriptor type
func (v *ValueDescriptor) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Id            *string      `json:"id"`
		Created       int64        `json:"created"`
		Description   *string      `json:"description"`
		Modified      int64        `json:"modified"`
		Origin        int64        `json:"origin"`
		Name          *string      `json:"name"`
		Min           *interface{} `json:"min"`
		Max           *interface{} `json:"max"`
		DefaultValue  *interface{} `json:"defaultValue"`
		Type          *string      `json:"type"`
		UomLabel      *string      `json:"uomLabel"`
		Formatting    *string      `json:"formatting"`
		Labels        []string     `json:"labels"`
		MediaType     *string      `json:"mediaType"`
		FloatEncoding *string      `json:"floatEncoding"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Set the fields
	if a.Id != nil {
		v.Id = *a.Id
	}
	if a.Description != nil {
		v.Description = *a.Description
	}
	if a.Name != nil {
		v.Name = *a.Name
	}
	if a.Min != nil {
		v.Min = *a.Min
	}
	if a.Max != nil {
		v.Max = *a.Max
	}
	if a.DefaultValue != nil {
		v.DefaultValue = *a.DefaultValue
	}
	if a.Type != nil {
		v.Type = *a.Type
	}
	if a.UomLabel != nil {
		v.UomLabel = *a.UomLabel
	}
	if a.Formatting != nil {
		v.Formatting = *a.Formatting
	}
	if a.MediaType != nil {
		v.MediaType = *a.MediaType
	}
	if a.FloatEncoding != nil {
		v.FloatEncoding = *a.FloatEncoding
	}
	v.Created = a.Created
	v.Modified = a.Modified
	v.Origin = a.Origin
	v.Labels = a.Labels

	v.isValidated, err = v.Validate()
	return err
}

// Validate satisfies the Validator interface
func (v ValueDescriptor) Validate() (bool, error) {
	if !v.isValidated {
		if v.Formatting != "" {
			formatSpecifier := "%(\\d+\\$)?([-#+ 0,(\\<]*)?(\\d+)?(\\.\\d+)?([tT])?([a-zA-Z%])"
			match, err := regexp.MatchString(formatSpecifier, v.Formatting)
			if err != nil {
				return false, NewErrContractInvalid(fmt.Sprintf("error validating format string: %s", v.Formatting))
			}
			if !match {
				return false, NewErrContractInvalid(fmt.Sprintf("format is not a valid printf format: %s", v.Formatting))
			}
		}
		if v.Name == "" {
			return false, NewErrContractInvalid("name for value descriptor not specified")
		}
	}
	return true, nil
}

// String returns a JSON encoded string representation of the model
func (a ValueDescriptor) String() string {
	out, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// From creates a ValueDescriptor based on the information provided in the DeviceResource.
func From(dr DeviceResource) ValueDescriptor {
	value := dr.Properties.Value
	units := dr.Properties.Units
	desc := ValueDescriptor{
		Name:          dr.Name,
		Min:           value.Minimum,
		Max:           value.Maximum,
		Type:          value.Type,
		UomLabel:      units.DefaultValue,
		DefaultValue:  value.DefaultValue,
		Formatting:    defaultValueDescriptorFormat,
		Description:   dr.Description,
		FloatEncoding: value.FloatEncoding,
		MediaType:     value.MediaType,
	}

	return desc
}
