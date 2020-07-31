// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

type serverConfig struct {
	port string `json:"port"`
}

type chanConfig struct {
	control string `json:"control"`
	data    string `json:"data"`
}

type edgexConfig struct {
	url string `json:"url"`
}

type logConfig struct {
	level string `json:"level"`
}

type mqttConfig struct {
	url         string `json:"url"`
	username    string `json:"username"`
	password    string `json:"json"`
	mtls        bool   `json:"mtls"`
	skipTLSVer  bool   `json:"skip_tls_ver"`
	retain      bool   `json:"retain"`
	QoS         byte   `json:"qos"`
	caPath      string `json:"ca_path"`
	certPath    string `json:"cert_path"`
	privKeyPath string `json:"priv_key_path"`
}

// Config struct of Mainflux Agent
type agentConfig struct {
	server   serverConfig `json:"server"`
	channels chanConfig   `json:"channels"`
	edgex    edgexConfig  `json:"edgex"`
	log      logConfig    `json:"log"`
	mqtt     mqttConfig   `json:"mqtt"`
}
