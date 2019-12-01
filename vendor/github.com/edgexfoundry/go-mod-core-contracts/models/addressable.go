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
	"strconv"
	"strings"
)

// Addressable holds information indicating how to contact a specific endpoint
type Addressable struct {
	Timestamps
	Id          string `json:"id,omitempty"`          // ID is a unique identifier for the Addressable, such as a UUID
	Name        string `json:"name,omitempty"`        // Name is a unique name given to the Addressable
	Protocol    string `json:"protocol,omitempty"`    // Protocol for the address (HTTP/TCP)
	HTTPMethod  string `json:"method,omitempty"`      // Method for connecting (i.e. POST)
	Address     string `json:"address,omitempty"`     // Address of the addressable
	Port        int    `json:"port,omitempty,Number"` // Port for the address
	Path        string `json:"path,omitempty"`        // Path for callbacks
	Publisher   string `json:"publisher,omitempty"`   // For message bus protocols
	User        string `json:"user,omitempty"`        // User id for authentication
	Password    string `json:"password,omitempty"`    // Password of the user for authentication for the addressable
	Topic       string `json:"topic,omitempty"`       // Topic for message bus addressables
	isValidated bool   // internal member used for validation check
}

type addressableAlias Addressable

// MarshalJSON implements the Marshaler interface for the Addressable type
// Use custom logic to create the URL and Base URL
func (a Addressable) MarshalJSON() ([]byte, error) {
	aux := struct {
		addressableAlias
		BaseURL string `json:"baseURL,omitempty"`
		URL     string `json:"url,omitempty"`
	}{
		addressableAlias: addressableAlias(a),
	}

	if a.Protocol != "" && a.Address != "" {
		// Get the base URL
		aux.BaseURL = a.GetBaseURL()

		// Get the URL
		aux.URL = aux.BaseURL
		if a.Publisher == "" && a.Topic != "" {
			aux.URL += a.Topic + "/"
		}
		aux.URL += a.Path
	}

	return json.Marshal(aux)
}

// UnmarshalJSON implements the Unmarshaler interface for the Addressable type
func (a *Addressable) UnmarshalJSON(data []byte) error {
	var err error
	var alias addressableAlias
	if err = json.Unmarshal(data, &alias); err != nil {
		return err
	}

	*a = Addressable(alias)
	a.isValidated, err = a.Validate()

	return err
}

// Validate satisfies the Validator interface
func (a Addressable) Validate() (bool, error) {
	if !a.isValidated {
		if a.Id == "" && a.Name == "" {
			return false, NewErrContractInvalid("Addressable ID and Name are both blank")
		}
		return true, nil
	}
	return a.isValidated, nil
}

// String returns a JSON encoded string representation of the addressable.
func (a Addressable) String() string {
	out, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// GetBaseURL returns a base URL consisting of protocol, host and port as a string assembled from the constituent parts of the Addressable
func (a Addressable) GetBaseURL() string {
	protocol := strings.ToLower(a.Protocol)
	address := a.Address
	port := strconv.Itoa(a.Port)
	baseUrl := protocol + "://" + address + ":" + port
	return baseUrl
}

// GetCallbackURL returns the callback url for the addressable if all relevant tokens have values.
// If any token is missing, string will be empty. Tokens include protocol, address, port and path.
func (a Addressable) GetCallbackURL() string {
	url := ""
	if len(a.Protocol) > 0 && len(a.Address) > 0 && a.Port > 0 && len(a.Path) > 0 {
		url = a.GetBaseURL() + a.Path
	}

	return url
}
