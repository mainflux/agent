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

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
)

// CommandResponse identifies a specific device along with its supported commands.
type CommandResponse struct {
	Id             string         `json:"id,omitempty"`             // Id uniquely identifies the CommandResponse, UUID for example.
	Name           string         `json:"name,omitempty"`           // Unique name for identifying a device
	AdminState     AdminState     `json:"adminState,omitempty"`     // Admin state (locked/unlocked)
	OperatingState OperatingState `json:"operatingState,omitempty"` // Operating state (enabled/disabled)
	LastConnected  int64          `json:"lastConnected,omitempty"`  // Time (milliseconds) that the device last provided any feedback or responded to any request
	LastReported   int64          `json:"lastReported,omitempty"`   // Time (milliseconds) that the device reported data to the core microservice
	Labels         []string       `json:"labels,omitempty"`         // Other labels applied to the device to help with searching
	Location       interface{}    `json:"location,omitempty"`       // Device service specific location (interface{} is an empty interface so it can be anything)
	Commands       []Command      `json:"commands,omitempty"`       // Associated Device Profile - Describes the device
}

/*
 * String function for representing a device
 */
func (d CommandResponse) String() string {
	out, err := json.Marshal(d)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

/*
 * CommandResponseFromDevice will create a CommandResponse struct from the supplied Device struct
 */
func CommandResponseFromDevice(d Device, commands []Command, cmdURL string) CommandResponse {
	cmdResp := CommandResponse{
		Id:             d.Id,
		Name:           d.Name,
		AdminState:     d.AdminState,
		OperatingState: d.OperatingState,
		LastConnected:  d.LastConnected,
		LastReported:   d.LastReported,
		Labels:         d.Labels,
		Location:       d.Location,
		Commands:       commands,
	}

	basePath := fmt.Sprintf("%s%s/%s/command/", cmdURL, clients.ApiDeviceRoute, d.Id)
	// TODO: Find a way to encapsulate this within the "Action" struct if possible
	for i := 0; i < len(cmdResp.Commands); i++ {
		url := basePath + cmdResp.Commands[i].Id
		cmdResp.Commands[i].Get.URL = url
		cmdResp.Commands[i].Put.URL = url
	}

	return cmdResp
}
