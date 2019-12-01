/*******************************************************************************
 * Copyright 2019 Dell Technologies Inc.
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
 *
 *******************************************************************************/

package models

import (
	"encoding/json"
	"reflect"
)

type Transmission struct {
	Timestamps
	ID           string               `json:"id"`
	Notification Notification         `json:"notification"`
	Receiver     string               `json:"receiver,omitempty"`
	Channel      Channel              `json:"channel,omitempty"`
	Status       TransmissionStatus   `json:"status,omitempty"`
	ResendCount  int                  `json:"resendcount"`
	Records      []TransmissionRecord `json:"records,omitempty"`
	isValidated  bool
}

// Marshal returns a JSON encoded byte array representation of the model
func (t Transmission) MarshalJSON() ([]byte, error) {
	alias := struct {
		Timestamps
		ID           string               `json:"id,omitempty"`
		Notification *Notification        `json:"notification,omitempty"`
		Receiver     string               `json:"receiver,omitempty"`
		Channel      *Channel             `json:"channel,omitempty"`
		Status       TransmissionStatus   `json:"status,omitempty"`
		ResendCount  *int                 `json:"resendcount,omitempty"`
		Records      []TransmissionRecord `json:"records,omitempty"`
	}{
		Timestamps:   t.Timestamps,
		ID:           t.ID,
		Notification: &t.Notification,
		Receiver:     t.Receiver,
		Channel:      &t.Channel,
		Status:       t.Status,
		ResendCount:  &t.ResendCount,
		Records:      t.Records,
	}

	// if we don't use omitempty, then an empty object always has a ResendCount included
	// if we do use omitempty, then a ResendCount of 0 is not included in the object when it should be
	if reflect.DeepEqual(t, Transmission{}) {
		alias.ResendCount = nil
	}
	// do not marshal empty member objects
	if reflect.DeepEqual(t.Notification, Notification{}) {
		alias.Notification = nil
	}
	if reflect.DeepEqual(t.Channel, Channel{}) {
		alias.Channel = nil
	}

	return json.Marshal(alias)
}

// UnmarshalJSON implements the Unmarshaler interface for the Transmission type
func (t *Transmission) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Timestamps
		ID           *string              `json:"id"`
		Notification Notification         `json:"notification,omitempty"`
		Receiver     *string              `json:"receiver,omitempty"`
		Channel      Channel              `json:"channel,omitempty"`
		Status       TransmissionStatus   `json:"status,omitempty"`
		ResendCount  int                  `json:"resendcount"`
		Records      []TransmissionRecord `json:"records,omitempty"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}
	// Nillable fields
	if a.ID != nil {
		t.ID = *a.ID
	}
	if a.Receiver != nil {
		t.Receiver = *a.Receiver
	}

	t.Notification = a.Notification
	t.Channel = a.Channel
	t.Status = a.Status
	t.ResendCount = a.ResendCount
	t.Records = a.Records
	t.Timestamps = a.Timestamps

	t.isValidated, err = t.Validate()

	return err
}

// Validate satisfies the Validator interface
func (t Transmission) Validate() (bool, error) {
	if !t.isValidated {

		if t.Notification.Slug == "" {
			return false, NewErrContractInvalid("Transmission's Notification is blank")
		}
		if t.Receiver == "" {
			return false, NewErrContractInvalid("Transmission's Receiver is blank")
		}
		if t.Channel.Type == "" {
			return false, NewErrContractInvalid("Transmission's Channel is blank")
		}
		if t.Status == "" {
			return false, NewErrContractInvalid("Transmission's Status is blank")
		}
		if t.ResendCount < 0 {
			return false, NewErrContractInvalid("Transmission's ResendCount is blank")
		}

		err := validate(t)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return t.isValidated, nil
}

// String returns a JSON encoded string representation of the model
func (t Transmission) String() string {
	out, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
