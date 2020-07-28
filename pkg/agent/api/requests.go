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
	agent agentConfig
	file  string
}

func (req addConfigReq) validate() error {
	if req.agent.server.port == "" ||
		req.agent.mqtt.username == "" ||
		req.agent.mqtt.password == "" ||
		req.agent.channels.control == "" ||
		req.agent.channels.data == "" ||
		req.agent.log.level == "" ||
		req.agent.edgex.url == "" ||
		req.agent.mqtt.url == "" {
		return agent.ErrMalformedEntity
	}

	return nil
}
