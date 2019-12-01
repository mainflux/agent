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
	"strings"
)

// OperatingState Constant String
type OperatingState string

/*
	Enabled  : ENABLED
	Disabled : DISABLED
*/
const (
	Enabled  = "ENABLED"
	Disabled = "DISABLED"
)

// UnmarshalJSON : Struct into json
func (os *OperatingState) UnmarshalJSON(data []byte) error {
	// Extract the string from data.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("OperatingState should be a string, got %s", data)
	}

	new := OperatingState(strings.ToUpper(s))
	*os = new

	return nil
}

// Validate satisfies the Validator interface
func (os OperatingState) Validate() (bool, error) {
	_, found := map[string]OperatingState{"ENABLED": Enabled, "DISABLED": Disabled}[string(os)]
	if !found {
		return false, NewErrContractInvalid(fmt.Sprintf("invalid OperatingState %q", os))
	}
	return true, nil
}

// GetOperatingState is called from within the router logic of the core services. For example, there are PUT calls
// like the one below from core-metadata which specify their update parameters in the URL
//
// d.HandleFunc("/{"+ID+"}/"+OPSTATE+"/{"+OPSTATE+"}", restSetDeviceOpStateById).Methods(http.MethodPut)
//
// Updates like this should be refactored to pass a body containing the new values instead of via the URL. This
// would allow us to utilize the model validation above and remove the logic from the controller.
//
// This will be removed once work on the following issue begins -- https://github.com/edgexfoundry/edgex-go/issues/1244
func GetOperatingState(os string) (OperatingState, bool) {
	os = strings.ToUpper(os)
	o, err := map[string]OperatingState{"ENABLED": Enabled, "DISABLED": Disabled}[os]
	return o, err
}
