// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type ServerConf struct {
	Port    string `toml:"port"`
	NatsURL string `toml:"nats_url"`
}

type ThingConf struct {
	ID  string `toml:"id"`
	Key string `toml:"key"`
}

type ChanConf struct {
	Control string `toml:"control"`
	Data    string `toml:"data"`
}

type EdgexConf struct {
	URL string `toml:"url"`
}

type LogConf struct {
	Level string `toml:"level"`
}

type MQTTConf struct {
	URL string `toml:"url"`
}

// Config struct of Mainflux Agent
type AgentConf struct {
	Server   ServerConf `toml:"server"`
	Thing    ThingConf  `toml:"thing"`
	Channels ChanConf   `toml:"channels"`
	Edgex    EdgexConf  `toml:"edgex"`
	Log      LogConf    `toml:"log"`
	MQTT     MQTTConf   `toml:"mqtt"`
}

type Config struct {
	Agent AgentConf
	File  string
}

func New(sc ServerConf, tc ThingConf, cc ChanConf, ec EdgexConf, lc LogConf, mc MQTTConf, file string) *Config {
	ac := AgentConf{
		Server:   sc,
		Thing:    tc,
		Channels: cc,
		Edgex:    ec,
		Log:      lc,
		MQTT:     mc,
	}

	return &Config{
		Agent: ac,
		File:  file,
	}
}

// Save - store config in a file
func (c *Config) Save() error {
	b, err := toml.Marshal(*c)
	if err != nil {
		fmt.Printf("Error reading config file: %s", err)
		return err
	}

	if err := ioutil.WriteFile(c.File, b, 0644); err != nil {
		fmt.Printf("Error writing toml: %s", err)
		return err
	}

	return nil
}

// Read - retrieve config from a file
func (c *Config) Read() error {
	data, err := ioutil.ReadFile(c.File)
	if err != nil {
		fmt.Printf("Error reading config file: %s", err)
		return err
	}

	if err := toml.Unmarshal(data, c); err != nil {
		fmt.Printf("Error unmarshaling toml: %s", err)
		return err
	}

	return nil
}
