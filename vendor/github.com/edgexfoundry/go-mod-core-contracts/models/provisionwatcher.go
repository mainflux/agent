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

type ProvisionWatcher struct {
	Timestamps
	Id                  string              `json:"id"`
	Name                string              `json:"name"`                // unique name and identifier of the provision watcher
	Identifiers         map[string]string   `json:"identifiers"`         // set of key value pairs that identify property (MAC, HTTP,...) and value to watch for (00-05-1B-A1-99-99, 10.0.0.1,...)
	BlockingIdentifiers map[string][]string `json:"blockingidentifiers"` // set of key-values pairs that identify devices which will not be added despite matching on Identifiers
	Profile             DeviceProfile       `json:"profile"`             // device profile that should be applied to the devices available at the identifier addresses
	Service             DeviceService       `json:"service"`             // device service that new devices will be associated to
	AdminState          AdminState          `json:"adminState"`          // administrative state for new devices - either unlocked or locked
	OperatingState      OperatingState      `validate:"-"`               // Deprecated: exists for historical compatibility and will be ignored
	isValidated         bool                ``                           // internal member used for validation check
}

// MarshalJSON returns a JSON encoded byte representation of the model
func (pw ProvisionWatcher) MarshalJSON() ([]byte, error) {
	test := struct {
		Timestamps
		Id                  string               `json:"id,omitempty"`
		Name                string               `json:"name,omitempty"`                // unique name and identifier of the addressable
		Identifiers         *map[string]string   `json:"identifiers,omitempty"`         // set of key value pairs that identify property (MAC, HTTP,...) and value to watch for (00-05-1B-A1-99-99, 10.0.0.1,...)
		BlockingIdentifiers *map[string][]string `json:"blockingidentifiers,omitempty"` // set of key-values pairs that identify devices which will not be added despite matching on Identifiers
		Profile             *DeviceProfile       `json:"profile,omitempty"`             // device profile that should be applied to the devices available at the identifier addresses
		Service             *DeviceService       `json:"service,omitempty"`             // device service that new devices will be associated to
		AdminState          AdminState           `json:"adminState,omitempty"`          // administrative state for new devices - either unlocked or locked
	}{
		Timestamps:          pw.Timestamps,
		Id:                  pw.Id,
		Name:                pw.Name,
		Identifiers:         &pw.Identifiers,
		BlockingIdentifiers: &pw.BlockingIdentifiers,
		Profile:             &pw.Profile,
		Service:             &pw.Service,
		AdminState:          pw.AdminState,
	}

	// Empty maps are null
	if len(pw.Identifiers) == 0 {
		test.Identifiers = nil
	}
	if len(pw.BlockingIdentifiers) == 0 {
		test.BlockingIdentifiers = nil
	}

	// Empty objects are nil
	if reflect.DeepEqual(pw.Profile, DeviceProfile{}) {
		test.Profile = nil
	}
	if reflect.DeepEqual(pw.Service, DeviceService{}) {
		test.Service = nil
	}

	return json.Marshal(test)
}

// UnmarshalJSON implements the Unmarshaler interface for the ProvisionWatcher type
func (pw *ProvisionWatcher) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Timestamps          `json:",inline"`
		Id                  string              `json:"id"`
		Name                *string             `json:"name"`
		Identifiers         map[string]string   `json:"identifiers"`
		BlockingIdentifiers map[string][]string `json:"blockingidentifiers"`
		Profile             DeviceProfile       `json:"profile"`
		Service             DeviceService       `json:"service"`
		AdminState          AdminState          `json:"adminState"`
	}
	a := Alias{}

	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Name can be nil
	if a.Name != nil {
		pw.Name = *a.Name
	}
	pw.Timestamps = a.Timestamps
	pw.Id = a.Id
	pw.Identifiers = a.Identifiers
	pw.BlockingIdentifiers = a.BlockingIdentifiers
	pw.Profile = a.Profile
	pw.Service = a.Service
	pw.AdminState = a.AdminState

	pw.isValidated, err = pw.Validate()

	return err
}

// Validate satisfies the Validator interface
func (pw ProvisionWatcher) Validate() (bool, error) {
	if !pw.isValidated {
		if pw.Name == "" {
			return false, NewErrContractInvalid("provision watcher name is blank")
		}
		err := validate(pw)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return pw.isValidated, nil
}

// String returns a JSON encoded string representation of the model
func (pw ProvisionWatcher) String() string {
	out, err := json.Marshal(pw)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
