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
)

// Channel supports transmissions and notifications with fields for delivery via email or REST
type Channel struct {
	Type          ChannelType `json:"type,omitempty"`          // Type indicates whether the channel facilitates email or REST
	MailAddresses []string    `json:"mailAddresses,omitempty"` // MailAddresses contains email addresses
	Url           string      `json:"url,omitempty"`           // URL contains a REST API destination
}

func (c Channel) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
