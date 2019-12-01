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

type IntervalAction struct {
	ID          string `json:"id,omitempty"`
	Created     int64  `json:"created,omitempty"`
	Modified    int64  `json:"modified,omitempty"`
	Origin      int64  `json:"origin,omitempty"`
	Name        string `json:"name,omitempty"`
	Interval    string `json:"interval,omitempty"`
	Parameters  string `json:"parameters,omitempty"`
	Target      string `json:"target,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	HTTPMethod  string `json:"httpMethod,omitempty"`
	Address     string `json:"address,omitempty"`
	Port        int    `json:"port,omitempty"`
	Path        string `json:"path,omitempty"`
	Publisher   string `json:"publisher,omitempty"`
	User        string `json:"user,omitempty"`
	Password    string `json:"password,omitempty"`
	Topic       string `json:"topic,omitempty"`
	isValidated bool   // internal member used for validation check
}

// UnmarshalJSON implements the Unmarshaler interface for the IntervalAction type
func (ia *IntervalAction) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		ID         *string `json:"id"`
		Created    int64   `json:"created"`
		Modified   int64   `json:"modified"`
		Origin     int64   `json:"origin"`
		Name       *string `json:"name"`
		Interval   *string `json:"interval"`
		Parameters *string `json:"parameters"`
		Target     *string `json:"target"`
		Protocol   *string `json:"protocol"`
		HTTPMethod *string `json:"httpMethod"`
		Address    *string `json:"address"`
		Port       int     `json:"port"`
		Path       *string `json:"path"`
		Publisher  *string `json:"publisher"`
		User       *string `json:"user"`
		Password   *string `json:"password"`
		Topic      *string `json:"topic"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Nillable fields
	if a.ID != nil {
		ia.ID = *a.ID
	}
	if a.Name != nil {
		ia.Name = *a.Name
	}
	if a.Interval != nil {
		ia.Interval = *a.Interval
	}
	if a.Parameters != nil {
		ia.Parameters = *a.Parameters
	}
	if a.Target != nil {
		ia.Target = *a.Target
	}
	if a.Protocol != nil {
		ia.Protocol = *a.Protocol
	}
	if a.HTTPMethod != nil {
		ia.HTTPMethod = *a.HTTPMethod
	}
	if a.Address != nil {
		ia.Address = *a.Address
	}
	if a.Path != nil {
		ia.Path = *a.Path
	}
	if a.Publisher != nil {
		ia.Publisher = *a.Publisher
	}
	if a.User != nil {
		ia.User = *a.User
	}
	if a.Password != nil {
		ia.Password = *a.Password
	}
	if a.Topic != nil {
		ia.Topic = *a.Topic
	}
	ia.Created = a.Created
	ia.Modified = a.Modified
	ia.Origin = a.Origin
	ia.Port = a.Port

	ia.isValidated, err = ia.Validate()

	return err
}

// Validate satisfies the Validator interface
func (ia IntervalAction) Validate() (bool, error) {
	if !ia.isValidated {
		if ia.ID == "" && ia.Name == "" {
			return false, NewErrContractInvalid("IntervalAction ID and Name are both blank")
		}
		if ia.Target == "" {
			return false, NewErrContractInvalid("intervalAction target is blank")
		}
		if ia.Interval == "" {
			return false, NewErrContractInvalid("intervalAction interval is blank")
		}
		return true, nil
	}
	return ia.isValidated, nil
}

func (ia IntervalAction) String() string {
	out, err := json.Marshal(ia)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func (ia IntervalAction) GetBaseURL() string {
	protocol := strings.ToLower(ia.Protocol)
	address := ia.Address
	port := strconv.Itoa(ia.Port)
	baseUrl := protocol + "://" + address + ":" + port
	return baseUrl
}
