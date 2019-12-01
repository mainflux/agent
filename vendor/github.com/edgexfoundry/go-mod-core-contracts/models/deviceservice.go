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

// DeviceService represents a service that is responsible for proxying connectivity between a set of devices and the
// EdgeX Foundry core services.
type DeviceService struct {
	DescribedObject
	Id             string         `json:"id"`
	Name           string         `json:"name"`           // time in milliseconds that the device last provided any feedback or responded to any request
	LastConnected  int64          `json:"lastConnected"`  // time in milliseconds that the device last reported data to the core
	LastReported   int64          `json:"lastReported"`   // operational state - either enabled or disabled
	OperatingState OperatingState `json:"operatingState"` // operational state - ether enabled or disableddc
	Labels         []string       `json:"labels"`         // tags or other labels applied to the device service for search or other identification needs
	Addressable    Addressable    `json:"addressable"`    // address (MQTT topic, HTTP address, serial bus, etc.) for reaching the service
	AdminState     AdminState     `json:"adminState"`     // Device Service Admin State
	isValidated    bool           // internal member used for validation check
}

// MarshalJSON implements the Marshaler interface in order to make empty strings null
func (ds DeviceService) MarshalJSON() ([]byte, error) {
	test := struct {
		DescribedObject `json:",inline"`
		Id              string         `json:"id,omitempty"`
		Name            string         `json:"name,omitempty"`           // time in milliseconds that the device last provided any feedback or responded to any request
		LastConnected   int64          `json:"lastConnected,omitempty"`  // time in milliseconds that the device last reported data to the core
		LastReported    int64          `json:"lastReported,omitempty"`   // operational state - either enabled or disabled
		OperatingState  OperatingState `json:"operatingState,omitempty"` // operational state - ether enabled or disableddc
		Labels          []string       `json:"labels,omitempty"`         // tags or other labels applied to the device service for search or other identification needs
		Addressable     *Addressable   `json:"addressable,omitempty"`    // address (MQTT topic, HTTP address, serial bus, etc.) for reaching the service
		AdminState      AdminState     `json:"adminState,omitempty"`     // Device Service Admin State
	}{
		DescribedObject: ds.DescribedObject,
		Id:              ds.Id,
		Name:            ds.Name,
		LastConnected:   ds.LastConnected,
		LastReported:    ds.LastReported,
		OperatingState:  ds.OperatingState,
		Labels:          ds.Labels,
		Addressable:     &ds.Addressable,
		AdminState:      ds.AdminState,
	}

	if reflect.DeepEqual(*test.Addressable, Addressable{}) {
		test.Addressable = nil
	}

	return json.Marshal(test)
}

// UnmarshalJSON implements the Unmarshaler interface for the DeviceService type
func (ds *DeviceService) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		DescribedObject `json:",inline"`
		Id              string         `json:"id"`
		Name            *string        `json:"name"`           // time in milliseconds that the device last provided any feedback or responded to any request
		LastConnected   int64          `json:"lastConnected"`  // time in milliseconds that the device last reported data to the core
		LastReported    int64          `json:"lastReported"`   // operational state - either enabled or disabled
		OperatingState  OperatingState `json:"operatingState"` // operational state - ether enabled or disableddc
		Labels          []string       `json:"labels"`         // tags or other labels applied to the device service for search or other identification needs
		Addressable     Addressable    `json:"addressable"`    // address (MQTT topic, HTTP address, serial bus, etc.) for reaching the service
		AdminState      AdminState     `json:"adminState"`     // Device Service Admin State
	}
	a := Alias{}

	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Set the fields
	ds.AdminState = a.AdminState
	ds.DescribedObject = a.DescribedObject
	ds.LastConnected = a.LastConnected
	ds.LastReported = a.LastReported
	ds.OperatingState = a.OperatingState
	ds.Labels = a.Labels
	ds.Addressable = a.Addressable
	ds.Id = a.Id

	// Name can be nil
	if a.Name != nil {
		ds.Name = *a.Name
	}

	ds.isValidated, err = ds.Validate()

	return err
}

// Validate satisfies the Validator interface
func (ds DeviceService) Validate() (bool, error) {
	if !ds.isValidated {
		if ds.Id == "" && ds.Name == "" {
			return false, NewErrContractInvalid("Device Service ID and Name are both blank")
		}
		return true, nil
	}
	return ds.isValidated, nil
}

/*
 * To String function for DeviceService
 */
func (ds DeviceService) String() string {
	out, err := json.Marshal(ds)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
