// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

type serverConf struct {
	port string `json:"port"`
}

type chanConf struct {
	control string `json:"control"`
	data    string `json:"data"`
}

type edgexConf struct {
	url string `json:"url"`
}

type logConf struct {
	level string `json:"level"`
}

type mqttConf struct {
	url         string `json:"url"`
	username    string `json:"username"`
	password    string `json:"json"`
	mtls        bool   `json:"mtls"`
	skipTLSVer  bool   `json:"skip_tls_ver"`
	retain      bool   `json:"retain"`
	qoS         int    `json:"qos"`
	caPath      string `json:"ca_path"`
	certPath    string `json:"cert_path"`
	privKeyPath string `json:"priv_key_path"`
}

// Config struct of Mainflux Agent
type agentConf struct {
	server   serverConf `json:"server"`
	channels chanConf   `json:"channels"`
	edgex    edgexConf  `json:"edgex"`
	log      logConf    `json:"log"`
	mqtt     mqttConf   `json:"mqtt"`
}
