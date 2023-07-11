// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

type serverConfig struct {
	Port string `json:"port"`
}

type chanConfig struct {
	Control string `json:"control"`
	Data    string `json:"data"`
}

type edgexConfig struct {
	Url string `json:"url"`
}

type logConfig struct {
	Level string `json:"level"`
}

type mqttConfig struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"json"`
	QoS      byte   `json:"qos"`
}

// Config struct of Mainflux Agent
type agentConfig struct {
	Server   serverConfig `json:"server"`
	Channels chanConfig   `json:"channels"`
	Edgex    edgexConfig  `json:"edgex"`
	Log      logConfig    `json:"log"`
	Mqtt     mqttConfig   `json:"mqtt"`
}
