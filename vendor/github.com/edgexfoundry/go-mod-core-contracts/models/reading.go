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
	"strings"
)

// Constants related to Reading ValueTypes
const (
	ValueTypeBool         = "Bool"
	ValueTypeString       = "String"
	ValueTypeUint8        = "Uint8"
	ValueTypeUint16       = "Uint16"
	ValueTypeUint32       = "Uint32"
	ValueTypeUint64       = "Uint64"
	ValueTypeInt8         = "Int8"
	ValueTypeInt16        = "Int16"
	ValueTypeInt32        = "Int32"
	ValueTypeInt64        = "Int64"
	ValueTypeFloat32      = "Float32"
	ValueTypeFloat64      = "Float64"
	ValueTypeBinary       = "Binary"
	ValueTypeBoolArray    = "BoolArray"
	ValueTypeStringArray  = "StringArray"
	ValueTypeUint8Array   = "Uint8Array"
	ValueTypeUint16Array  = "Uint16Array"
	ValueTypeUint32Array  = "Uint32Array"
	ValueTypeUint64Array  = "Uint64Array"
	ValueTypeInt8Array    = "Int8Array"
	ValueTypeInt16Array   = "Int16Array"
	ValueTypeInt32Array   = "Int32Array"
	ValueTypeInt64Array   = "Int64Array"
	ValueTypeFloat32Array = "Float32Array"
	ValueTypeFloat64Array = "Float64Array"
)

// Reading contains data that was gathered from a device.
//
// NOTE a Reading's BinaryValue is not to be persisted in the database. This architectural decision requires that
// serialization validation be relaxed for enforcing the presence of binary data for Binary ValueTypes. Also, that
// issuing GET operations to obtain Readings directly or indirectly via Events will result in a Reading with no
// BinaryValue for Readings with a ValueType of Binary. BinaryValue is to be present when creating or updating a Reading
// either directly, indirectly via an Event, and when the information is put on the EventBus.
type Reading struct {
	Id            string `json:"id,omitempty" codec:"id,omitempty"`
	Pushed        int64  `json:"pushed,omitempty" codec:"pushed,omitempty"`   // When the data was pushed out of EdgeX (0 - not pushed yet)
	Created       int64  `json:"created,omitempty" codec:"created,omitempty"` // When the reading was created
	Origin        int64  `json:"origin,omitempty" codec:"origin,omitempty"`
	Modified      int64  `json:"modified,omitempty" codec:"modified,omitempty"`
	Device        string `json:"device,omitempty" codec:"device,omitempty"`
	Name          string `json:"name,omitempty" codec:"name,omitempty"`
	Value         string `json:"value,omitempty" codec:"value,omitempty"` // Device sensor data value
	ValueType     string `json:"valueType,omitempty" codec:"valueType,omitempty"`
	FloatEncoding string `json:"floatEncoding,omitempty" codec:"floatEncoding,omitempty"`
	// BinaryValue binary data payload. This information is not persisted in the Database and is expected to be empty
	// when retrieving a Reading for the ValueType of Binary.
	BinaryValue []byte `json:"binaryValue,omitempty" codec:"binaryValue,omitempty"`
	MediaType   string `json:"mediaType,omitempty" codec:"mediaType,omitempty"`
	isValidated bool   // internal member used for validation check
}

// UnmarshalJSON implements the Unmarshaler interface for the Reading type
func (r *Reading) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Id            *string `json:"id"`
		Pushed        int64   `json:"pushed"`
		Created       int64   `json:"created"`
		Origin        int64   `json:"origin"`
		Modified      int64   `json:"modified"`
		Device        *string `json:"device"`
		Name          *string `json:"name"`
		Value         *string `json:"value"`
		ValueType     *string `json:"valueType"`
		FloatEncoding *string `json:"floatEncoding"`
		BinaryValue   []byte  `json:"binaryValue"`
		MediaType     *string `json:"mediaType"`
	}
	a := Alias{}

	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Set the fields
	if a.Id != nil {
		r.Id = *a.Id
	}
	if a.Device != nil {
		r.Device = *a.Device
	}
	if a.Name != nil {
		r.Name = *a.Name
	}
	if a.Value != nil {
		r.Value = *a.Value
	}
	if a.ValueType != nil {
		r.ValueType = normalizeValueTypeCase(*a.ValueType)
	}
	if a.FloatEncoding != nil {
		r.FloatEncoding = *a.FloatEncoding
	}
	if a.MediaType != nil {
		r.MediaType = *a.MediaType
	}
	r.Pushed = a.Pushed
	r.Created = a.Created
	r.Origin = a.Origin
	r.Modified = a.Modified
	r.BinaryValue = a.BinaryValue

	r.isValidated, err = r.Validate()
	return err
}

// Validate satisfies the Validator interface
func (r Reading) Validate() (bool, error) {
	// Shortcut if Reading has already been validated
	if r.isValidated {
		return true, nil
	}

	if r.Name == "" {
		return false, NewErrContractInvalid("name for reading's value descriptor not specified")
	}
	// We do not expect the BinaryValue to always be present. This is due to an architectural decision to not persist
	// Binary readings to save on memory. This means that the BinaryValue is only expected to be populated when creating
	// a new reading or event. Otherwise the value will be empty as it will be coming from the database where we are
	// explicitly not storing the information.
	if r.ValueType != ValueTypeBinary && r.Value == "" {
		return false, NewErrContractInvalid("reading has no value")
	}

	// Even though we do not want to enforce the BinaryValue always being present for Readings, we still want to enforce
	// the MediaType being specified when the BinaryValue is provided. This will most likely only take affect when
	// creating and updating events or readings.
	if len(r.BinaryValue) != 0 && len(r.MediaType) == 0 {
		return false, NewErrContractInvalid("media type must be specified for binary values")
	}

	if (r.ValueType == ValueTypeFloat32 || r.ValueType == ValueTypeFloat64) && len(r.FloatEncoding) == 0 {
		return false, NewErrContractInvalid("float encoding must be specified for float values")
	}
	return true, nil
}

// normalizeValueTypeCase normalize the reading's valueType to upper camel case
func normalizeValueTypeCase(valueType string) string {
	normalized := strings.Title(strings.ToLower(valueType))
	normalized = strings.ReplaceAll(normalized, "array", "Array")
	return normalized
}

// String returns a JSON encoded string representation of the model
func (r Reading) String() string {
	out, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}

	return string(out)
}
