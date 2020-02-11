// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package writers contain the domaSavein concept definitions needed to
// support Mainflux writer services functionality.
package config

import (
	"crypto/tls"
	"io/ioutil"

	"github.com/mainflux/mainflux/errors"
	"github.com/pelletier/go-toml"
)

const (
	dfltFile = "config.toml"
)

var (
	errReadConfigFile         = errors.New("Error reading config file")
	errWritingConfigFile      = errors.New("Error writing config file")
	errUnmarshalConfigContent = errors.New("Error unmarshaling config file conent")
)

type MQTTConf struct {
	Host        string          `json:"host" toml:"host" mapstructure:"host"`
	Username    string          `json:"username" toml:"username" mapstructure:"username"`
	Password    string          `json:"password" toml:"password" mapstructure:"password"`
	MTLS        bool            `json:"mtls" toml:"mtls" mapstructure:"mtls"`
	SkipTLSVer  bool            `json:"skip_tls_ver" toml:"skip_tls_ver" mapstructure:"skip_tls_ver"`
	Retain      bool            `json:"retain" toml:"retain" mapstructure:"retain"`
	QoS         int             `json:"qos" toml:"qos" mapstructure:"qos"`
	Channel     string          `json:"channel" toml:"channel" mapstructure:"channel"`
	CAPath      string          `json:"ca_path" toml:"ca_path" mapstructure:"ca_path"`
	CertPath    string          `json:"cert_path" toml:"cert_path" mapstructure:"cert_path"`
	PrivKeyPath string          `json:"priv_key_path" toml:"priv_key_path" mapstructure:"priv_key_path"`
	CA          []byte          `json:"-" toml:"-"`
	Cert        tls.Certificate `json:"-" toml:"-"`
}

type ServerConf struct {
	NatsURL  string `json:"nats" toml:"nats" mapstructure:"nats"`
	LogLevel string `json:"log_level" toml:"log_level" mapstructure:"log_level"`
	Port     string `json:"port" toml:"port" mapstructure:"port"`
}

type Config struct {
	Server ServerConf `json:"exp" toml:"exp" mapstructure:"exp"`
	Routes []Route    `json:"routes" toml:"routes" mapstructure:"routes"`
	MQTT   MQTTConf   `json:"mqtt" toml:"mqtt" mapstructure:"mqtt"`
	File   string     `json:"-"`
}

type Route struct {
	MqttTopic string `json:"mqtt_topic" toml:"mqtt_topic" mapstructure:"mqtt_topic"`
	NatsTopic string `json:"nats_topic" toml:"nats_topic" mapstructure:"nats_topic"`
	SubTopic  string `json:"subtopic" toml:"subtopic" mapstructure:"subtopic"`
	Type      string `json:"type" toml:"type" mapstructure:"type"`
}

func NewConfig(sc ServerConf, rc []Route, mc MQTTConf, file string) *Config {
	ac := Config{
		Server: sc,
		Routes: rc,
		MQTT:   mc,
		File:   file,
	}

	return &ac
}

// Save - store config in a file
func (c *Config) Save() errors.Error {
	b, err := toml.Marshal(*c)
	if err != nil {
		return errors.Wrap(errReadConfigFile, err)
	}
	file := dfltFile
	if c.File != "" {
		file = c.File
	}
	if err := ioutil.WriteFile(file, b, 0644); err != nil {
		return errors.Wrap(errWritingConfigFile, err)
	}

	return nil
}

// ReadFile - retrieve config from a file
func (c *Config) ReadFile() errors.Error {
	file := c.File
	if file == "" {
		file = dfltFile
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(errReadConfigFile, err)
	}
	if err := toml.Unmarshal(data, c); err != nil {
		return errors.Wrap(errUnmarshalConfigContent, err)
	}
	c.File = file
	return nil
}

// ReadBytes - retrieve config from a byte
func (c *Config) ReadBytes(data []byte) errors.Error {
	if err := toml.Unmarshal(data, c); err != nil {
		return errors.Wrap(errUnmarshalConfigContent, err)
	}
	return nil
}
