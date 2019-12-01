// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

type serverConf struct {
	port string `json:"port"`
}

type thingConf struct {
	id  string `json:"id"`
	key string `json:"key"`
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
	url string `json:"url"`
}

// Config struct of Mainflux Agent
type agentConf struct {
	server   serverConf `json:"server"`
	thing    thingConf  `json:"thing"`
	channels chanConf   `json:"channels"`
	edgex    edgexConf  `json:"edgex"`
	log      logConf    `json:"log"`
	mqtt     mqttConf   `json:"mqtt"`
}
