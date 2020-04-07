// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/mainflux/mainflux/errors"
	"github.com/pelletier/go-toml"
)

type ServerConf struct {
	Port    string `toml:"port"`
	NatsURL string `toml:"nats_url"`
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
	URL         string          `json:"url" toml:"url"`
	Username    string          `json:"username" toml:"username" mapstructure:"username"`
	Password    string          `json:"password" toml:"password" mapstructure:"password"`
	MTLS        bool            `json:"mtls" toml:"mtls" mapstructure:"mtls"`
	SkipTLSVer  bool            `json:"skip_tls_ver" toml:"skip_tls_ver" mapstructure:"skip_tls_ver"`
	Retain      bool            `json:"retain" toml:"retain" mapstructure:"retain"`
	QoS         byte            `json:"qos" toml:"qos" mapstructure:"qos"`
	CAPath      string          `json:"ca_path" toml:"ca_path" mapstructure:"ca_path"`
	CertPath    string          `json:"cert_path" toml:"cert_path" mapstructure:"cert_path"`
	PrivKeyPath string          `json:"priv_key_path" toml:"priv_key_path" mapstructure:"priv_key_path"`
	CA          []byte          `json:"-" toml:"-"`
	Cert        tls.Certificate `json:"-" toml:"-"`
}

// Config struct of Mainflux Agent
type AgentConf struct {
	Server   ServerConf `toml:"server"`
	Channels ChanConf   `toml:"channels"`
	Edgex    EdgexConf  `toml:"edgex"`
	Log      LogConf    `toml:"log"`
	MQTT     MQTTConf   `toml:"mqtt"`
}

type Config struct {
	Agent AgentConf
	File  string
}

func New(sc ServerConf, cc ChanConf, ec EdgexConf, lc LogConf, mc MQTTConf, file string) Config {
	ac := AgentConf{
		Server:   sc,
		Channels: cc,
		Edgex:    ec,
		Log:      lc,
		MQTT:     mc,
	}

	return Config{
		Agent: ac,
		File:  file,
	}
}

// Save - store config in a file
func Save(c Config) error {
	b, err := toml.Marshal(c)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading config file: %s", err))
	}
	if err := ioutil.WriteFile(c.File, b, 0644); err != nil {
		return errors.New(fmt.Sprintf("Error writing toml: %s", err))
	}
	return nil
}

// Read - retrieve config from a file
func Read(file string) (Config, error) {
	data, err := ioutil.ReadFile(file)
	c := Config{}
	if err != nil {
		return c, errors.New(fmt.Sprintf("Error reading config file: %s", err))
	}

	if err := toml.Unmarshal(data, &c); err != nil {
		return Config{}, errors.New(fmt.Sprintf("Error unmarshaling toml: %s", err))
	}
	return c, nil
}
