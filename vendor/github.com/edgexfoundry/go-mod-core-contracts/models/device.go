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

// Device represents a registered device participating within the EdgeX Foundry ecosystem
type Device struct {
	DescribedObject
	Id             string                        `json:"id"`             // ID uniquely identifies the device, a UUID for example
	Name           string                        `json:"name"`           // Unique name for identifying a device
	AdminState     AdminState                    `json:"adminState"`     // Admin state (locked/unlocked)
	OperatingState OperatingState                `json:"operatingState"` // Operating state (enabled/disabled)
	Protocols      map[string]ProtocolProperties `json:"protocols"`      // A map of supported protocols for the given device
	LastConnected  int64                         `json:"lastConnected"`  // Time (milliseconds) that the device last provided any feedback or responded to any request
	LastReported   int64                         `json:"lastReported"`   // Time (milliseconds) that the device reported data to the core microservice
	Labels         []string                      `json:"labels"`         // Other labels applied to the device to help with searching
	Location       interface{}                   `json:"location"`       // Device service specific location (interface{} is an empty interface so it can be anything)
	Service        DeviceService                 `json:"service"`        // Associated Device Service - One per device
	Profile        DeviceProfile                 `json:"profile"`        // Associated Device Profile - Describes the device
	AutoEvents     []AutoEvent                   `json:"autoEvents"`     // A list of auto-generated events coming from the device
	isValidated    bool                          // internal member used for validation check
}

// ProtocolProperties contains the device connection information in key/value pair
type ProtocolProperties map[string]string

// MarshalJSON implements the Marshaler interface in order to make empty strings null
func (d Device) MarshalJSON() ([]byte, error) {
	test := struct {
		DescribedObject
		Id             string                        `json:"id,omitempty"`
		Name           string                        `json:"name,omitempty"`
		AdminState     AdminState                    `json:"adminState,omitempty"`
		OperatingState OperatingState                `json:"operatingState,omitempty"`
		Protocols      map[string]ProtocolProperties `json:"protocols,omitempty"`
		LastConnected  int64                         `json:"lastConnected,omitempty"`
		LastReported   int64                         `json:"lastReported,omitempty"`
		Labels         []string                      `json:"labels,omitempty"`
		Location       interface{}                   `json:"location,omitempty"`
		Service        *DeviceService                `json:"service,omitempty"`
		Profile        *DeviceProfile                `json:"profile,omitempty"`
		AutoEvents     []AutoEvent                   `json:"autoEvents,omitempty"`
	}{
		Id:              d.Id,
		Name:            d.Name,
		DescribedObject: d.DescribedObject,
		AdminState:      d.AdminState,
		OperatingState:  d.OperatingState,
		Protocols:       d.Protocols,
		LastConnected:   d.LastConnected,
		LastReported:    d.LastReported,
		Labels:          d.Labels,
		Location:        d.Location,
		Service:         &d.Service,
		Profile:         &d.Profile,
		AutoEvents:      d.AutoEvents,
	}

	if reflect.DeepEqual(*test.Service, DeviceService{}) {
		test.Service = nil
	}

	if reflect.DeepEqual(*test.Profile, DeviceProfile{}) {
		test.Profile = nil
	}

	return json.Marshal(test)
}

// UnmarshalJSON implements the Unmarshaler interface for the Device type
func (d *Device) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		DescribedObject `json:",inline"`
		Id              string                        `json:"id"`
		Name            string                        `json:"name"`
		AdminState      AdminState                    `json:"adminState"`
		OperatingState  OperatingState                `json:"operatingState"`
		Protocols       map[string]ProtocolProperties `json:"protocols"`
		LastConnected   int64                         `json:"lastConnected"`
		LastReported    int64                         `json:"lastReported"`
		Labels          []string                      `json:"labels"`
		Location        interface{}                   `json:"location"`
		Service         DeviceService                 `json:"service"`
		Profile         DeviceProfile                 `json:"profile"`
		AutoEvents      []AutoEvent                   `json:"autoEvents"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	d.Id = a.Id
	d.Name = a.Name
	d.DescribedObject = a.DescribedObject
	d.AdminState = a.AdminState
	d.OperatingState = a.OperatingState
	d.Protocols = a.Protocols
	d.LastConnected = a.LastConnected
	d.LastReported = a.LastReported
	d.Labels = a.Labels
	d.Location = a.Location
	d.Service = a.Service
	d.Profile = a.Profile
	d.AutoEvents = a.AutoEvents

	d.isValidated, err = d.Validate()

	return err
}

// Validate satisfies the Validator interface
func (d Device) Validate() (bool, error) {
	if !d.isValidated {
		if d.Id == "" && d.Name == "" {
			return false, NewErrContractInvalid("Device ID and Name are both blank")
		}
		if len(d.Protocols) == 0 {
			return false, NewErrContractInvalid("no supporting protocol specified for device")
		}
		err := validate(d)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return d.isValidated, nil
}

/*
 * String function for representing a device
 */
func (d Device) String() string {
	out, err := json.Marshal(d)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// AllAssociatedValueDescriptors returns all the associated value descriptors through Put command parameters and Put/Get command return values
func (d *Device) AllAssociatedValueDescriptors(vdNames *[]string) {
	// Get the value descriptors with a map (set) where the keys are the value descriptor names
	vdNamesMap := map[string]string{}
	for _, c := range d.Profile.CoreCommands {
		c.AllAssociatedValueDescriptors(&vdNamesMap)
	}

	// Add the map keys (value descriptor names) to the list
	for vd := range vdNamesMap {
		*vdNames = append(*vdNames, vd)
	}
}
