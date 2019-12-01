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

// Export destination types
const (
	DestMQTT        = "MQTT_TOPIC"
	DestZMQ         = "ZMQ_TOPIC"
	DestIotCoreMQTT = "IOTCORE_TOPIC"
	DestAzureMQTT   = "AZURE_TOPIC"
	DestRest        = "REST_ENDPOINT"
	DestXMPP        = "XMPP_TOPIC"
	DestAWSMQTT     = "AWS_TOPIC"
	DestInfluxDB    = "INFLUXDB_ENDPOINT"
)

// Compression algorithm types
const (
	CompNone = "NONE"
	CompGzip = "GZIP"
	CompZip  = "ZIP"
)

// Data format types
const (
	FormatJSON            = "JSON"
	FormatXML             = "XML"
	FormatSerialized      = "SERIALIZED"
	FormatIoTCoreJSON     = "IOTCORE_JSON"
	FormatAzureJSON       = "AZURE_JSON"
	FormatAWSJSON         = "AWS_JSON"
	FormatCSV             = "CSV"
	FormatThingsBoardJSON = "THINGSBOARD_JSON"
	FormatNOOP            = "NOOP"
)

const (
	NotifyUpdateAdd    = "add"
	NotifyUpdateUpdate = "update"
	NotifyUpdateDelete = "delete"
)

// Registration - Defines the registration details
// on the part of north side export clients
type Registration struct {
	ID          string            `json:"id,omitempty"`
	Created     int64             `json:"created"`
	Modified    int64             `json:"modified"`
	Origin      int64             `json:"origin"`
	Name        string            `json:"name,omitempty"`
	Addressable Addressable       `json:"addressable,omitempty"`
	Format      string            `json:"format,omitempty"`
	Filter      Filter            `json:"filter,omitempty"`
	Encryption  EncryptionDetails `json:"encryption,omitempty"`
	Compression string            `json:"compression,omitempty"`
	Enable      bool              `json:"enable"`
	Destination string            `json:"destination,omitempty"`
	isValidated bool              // internal member used for validation check
}

// UnmarshalJSON implements the Unmarshaler interface for the DeviceService type
func (r *Registration) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		ID          *string           `json:"id"`
		Created     int64             `json:"created"`
		Modified    int64             `json:"modified"`
		Origin      int64             `json:"origin"`
		Name        *string           `json:"name"`
		Addressable Addressable       `json:"addressable"`
		Format      *string           `json:"format"`
		Filter      Filter            `json:"filter"`
		Encryption  EncryptionDetails `json:"encryption"`
		Compression *string           `json:"compression"`
		Enable      bool              `json:"enable"`
		Destination *string           `json:"destination"`
	}
	a := Alias{}

	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Fields can be nil
	if a.ID != nil {
		r.ID = *a.ID
	}
	if a.Name != nil {
		r.Name = *a.Name
	}
	if a.Format != nil {
		r.Format = *a.Format
	}
	if a.Compression != nil {
		r.Compression = *a.Compression
	}
	if a.Destination != nil {
		r.Destination = *a.Destination
	}
	r.Created = a.Created
	r.Modified = a.Modified
	r.Origin = a.Origin
	r.Addressable = a.Addressable
	r.Filter = a.Filter
	r.Encryption = a.Encryption
	r.Enable = a.Enable

	r.isValidated, err = r.Validate()

	return err
}

// Validate satisfies the Validator interface
func (reg Registration) Validate() (bool, error) {
	if !reg.isValidated {
		if reg.Name == "" {
			return false, NewErrContractInvalid("Name is required")
		}

		if reg.Compression == "" {
			reg.Compression = CompNone
		}

		if reg.Compression != CompNone &&
			reg.Compression != CompGzip &&
			reg.Compression != CompZip {
			return false, NewErrContractInvalid(fmt.Sprintf("Compression invalid: %s", reg.Compression))
		}

		if reg.Format != FormatJSON &&
			reg.Format != FormatXML &&
			reg.Format != FormatSerialized &&
			reg.Format != FormatIoTCoreJSON &&
			reg.Format != FormatAzureJSON &&
			reg.Format != FormatAWSJSON &&
			reg.Format != FormatCSV &&
			reg.Format != FormatThingsBoardJSON &&
			reg.Format != FormatNOOP {
			return false, NewErrContractInvalid(fmt.Sprintf("Format invalid: %s", reg.Format))
		}

		if reg.Destination != DestMQTT &&
			reg.Destination != DestZMQ &&
			reg.Destination != DestIotCoreMQTT &&
			reg.Destination != DestAzureMQTT &&
			reg.Destination != DestAWSMQTT &&
			reg.Destination != DestRest &&
			reg.Destination != DestInfluxDB {
			return false, NewErrContractInvalid(fmt.Sprintf("Destination invalid: %s", reg.Destination))
		}

		if reg.Encryption.Algo == "" {
			reg.Encryption.Algo = EncNone
		}

		if reg.Encryption.Algo != EncNone &&
			reg.Encryption.Algo != EncAes {
			return false, NewErrContractInvalid(fmt.Sprintf("Encryption invalid: %s", reg.Encryption.Algo))
		}
		err := validate(reg)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return reg.isValidated, nil
}
