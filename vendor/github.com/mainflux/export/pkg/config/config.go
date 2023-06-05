// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package writers contain the domaSavein concept definitions needed to
// support Mainflux writer services functionality.
package config

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"

	"github.com/mainflux/mainflux/pkg/errors"
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

type MQTT struct {
	Host              string          `json:"host" toml:"host" mapstructure:"host"`
	Username          string          `json:"username" toml:"username" mapstructure:"username"`
	Password          string          `json:"password" toml:"password" mapstructure:"password"`
	MTLS              bool            `json:"mtls" toml:"mtls" mapstructure:"mtls"`
	SkipTLSVer        bool            `json:"skip_tls_ver" toml:"skip_tls_ver" mapstructure:"skip_tls_ver"`
	Retain            bool            `json:"retain" toml:"retain" mapstructure:"retain"`
	QoS               int             `json:"qos" toml:"qos" mapstructure:"qos"`
	CAPath            string          `json:"ca_path" toml:"ca_path" mapstructure:"ca_path"`
	ClientCertPath    string          `json:"client_cert_path" toml:"client_cert_path" mapstructure:"client_cert_path"`
	ClientPrivKeyPath string          `json:"client_priv_key_path" toml:"client_priv_key_path" mapstructure:"client_priv_key_path"`
	ClientCert        string          `json:"client_cert" toml:"client_cert" mapstructure:"client_cert"`
	ClientCertKey     string          `json:"client_cert_key" toml:"client_cert_key" mapstructure:"client_cert_key"`
	CA                []byte          `json:"-" toml:"-"`
	TLSCert           tls.Certificate `json:"-" toml:"-"`
}

type Server struct {
	NatsURL   string `json:"nats" toml:"nats" mapstructure:"nats"`
	LogLevel  string `json:"log_level" toml:"log_level" mapstructure:"log_level"`
	Port      string `json:"port" toml:"port" mapstructure:"port"`
	CacheURL  string `json:"cache_url" toml:"cache_url" mapstructure:"port"`
	CachePass string `json:"cache_pass" toml:"cache_pass" mapstructure:"port"`
	CacheDB   string `json:"cache_db" toml:"cache_db" mapstructure:"port"`
}

type Config struct {
	Server Server  `json:"exp" toml:"exp" mapstructure:"exp"`
	Routes []Route `json:"routes" toml:"routes" mapstructure:"routes"`
	MQTT   MQTT    `json:"mqtt" toml:"mqtt" mapstructure:"mqtt"`
	File   string  `json:"file"`
}

type Route struct {
	MqttTopic string `json:"mqtt_topic" toml:"mqtt_topic" mapstructure:"mqtt_topic"`
	NatsTopic string `json:"nats_topic" toml:"nats_topic" mapstructure:"nats_topic"`
	SubTopic  string `json:"subtopic" toml:"subtopic" mapstructure:"subtopic"`
	Type      string `json:"type" toml:"type" mapstructure:"type"`
	Workers   int    `json:"workers" toml:"workers" mapstructure:"workers"`
}

// Save - store config in a file
func Save(c Config) error {
	b, err := toml.Marshal(c)
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
func ReadFile(file string) (Config, error) {
	c := Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return c, errors.Wrap(errReadConfigFile, err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, errors.Wrap(errUnmarshalConfigContent, err)
	}
	c.File = file
	return c, nil
}

// ReadBytes - read config from a bytes
func ReadBytes(data []byte) (Config, error) {
	c := Config{}
	e := toml.Unmarshal(data, &c)
	if e == nil {
		return c, nil
	}
	err := errors.Wrap(errUnmarshalConfigContent, e)
	if e := json.Unmarshal(data, &c); e != nil {
		return c, errors.Wrap(err, e)
	}
	return c, nil
}
