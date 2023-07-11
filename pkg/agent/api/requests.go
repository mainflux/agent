// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/mainflux/agent/pkg/agent"
)

type pubReq struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

func (req pubReq) validate() error {
	if req.Topic == "" || req.Payload == "" {
		return agent.ErrMalformedEntity
	}

	return nil
}

type execReq struct {
	BaseName string `json:"bn"`
	Name     string `json:"n"`
	Value    string `json:"vs"`
}

func (req execReq) validate() error {
	if req.BaseName == "" || req.Name != "exec" || req.Value == "" {
		return agent.ErrMalformedEntity
	}

	return nil
}

type addConfigReq struct {
	Agent agentConfig
}

func (req addConfigReq) validate() error {
	if req.Agent.Server.Port == "" ||
		req.Agent.Mqtt.Username == "" ||
		req.Agent.Mqtt.Password == "" ||
		req.Agent.Channels.Control == "" ||
		req.Agent.Channels.Data == "" ||
		req.Agent.Log.Level == "" ||
		req.Agent.Edgex.Url == "" ||
		req.Agent.Mqtt.Url == "" {
		return agent.ErrMalformedEntity
	}

	return nil
}
