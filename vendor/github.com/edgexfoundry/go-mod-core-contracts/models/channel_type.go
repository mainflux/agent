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
	"fmt"
)

// ChannelType controls the range of values which constitute valid delivery types for channels
type ChannelType string

const (
	Rest  = "REST"
	Email = "EMAIL"
)

// UnmarshalJSON implements the Unmarshaler interface for the type
func (as *ChannelType) UnmarshalJSON(data []byte) error {
	// Extract the string from data.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ChannelType should be a string, got %s", data)
	}

	got, err := map[string]ChannelType{"REST": Rest, "EMAIL": Email}[s]
	if !err {
		return fmt.Errorf("invalid ChannelType %q", s)
	}
	*as = got
	return nil
}

func (as ChannelType) Validate() (bool, error) {
	_, err := map[string]ChannelType{"REST": Rest, "EMAIL": Email}[string(as)]
	if !err {
		return false, NewErrContractInvalid(fmt.Sprintf("invalid Channeltype %q", as))
	}
	return true, nil
}
