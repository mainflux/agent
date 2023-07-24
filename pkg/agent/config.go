// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/pelletier/go-toml"
)

type ServerConfig struct {
	Port    string `toml:"port" json:"port"`
	NatsURL string `toml:"nats_url" json:"nats_url"`
}

type ChanConfig struct {
	Control string `toml:"control"`
	Data    string `toml:"data"`
}

type EdgexConfig struct {
	URL string `toml:"url"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

type MQTTConfig struct {
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
	ClientCert  string          `json:"client_cert" toml:"client_cert"`
	ClientKey   string          `json:"client_key" toml:"client_key"`
	CaCert      string          `json:"ca_cert" toml:"ca_cert"`
}

type HeartbeatConfig struct {
	Interval time.Duration `toml:"interval"`
}

type TerminalConfig struct {
	SessionTimeout time.Duration `toml:"session_timeout" json:"session_timeout"`
}

type Config struct {
	Server    ServerConfig    `toml:"server" json:"server"`
	Terminal  TerminalConfig  `toml:"terminal" json:"terminal"`
	Heartbeat HeartbeatConfig `toml:"heartbeat" json:"heartbeat"`
	Channels  ChanConfig      `toml:"channels" json:"channels"`
	Edgex     EdgexConfig     `toml:"edgex" json:"edgex"`
	Log       LogConfig       `toml:"log" json:"log"`
	MQTT      MQTTConfig      `toml:"mqtt" json:"mqtt"`
	File      string
}

func NewConfig(sc ServerConfig, cc ChanConfig, ec EdgexConfig, lc LogConfig, mc MQTTConfig, hc HeartbeatConfig, tc TerminalConfig, file string) Config {
	return Config{
		Server:    sc,
		Channels:  cc,
		Edgex:     ec,
		Log:       lc,
		MQTT:      mc,
		Heartbeat: hc,
		Terminal:  tc,
		File:      file,
	}
}

// Save - store config in a file.
func SaveConfig(c Config) error {
	b, err := toml.Marshal(c)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading config file: %s", err))
	}
	if err := os.WriteFile(c.File, b, 0644); err != nil {
		return errors.New(fmt.Sprintf("Error writing toml: %s", err))
	}
	return nil
}

// Read - retrieve config from a file.
func ReadConfig(file string) (Config, error) {
	data, err := os.ReadFile(file)
	c := Config{}
	if err != nil {
		return c, errors.New(fmt.Sprintf("Error reading config file: %s", err))
	}

	if err := toml.Unmarshal(data, &c); err != nil {
		return Config{}, errors.New(fmt.Sprintf("Error unmarshaling toml: %s", err))
	}
	return c, nil
}

// UnmarshalJSON parses the duration from JSON.
func (d *HeartbeatConfig) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	interval, ok := v["interval"]
	if !ok {
		return errors.New("missing value")
	}
	switch value := interval.(type) {
	case float64:
		d.Interval = time.Duration(value)
		return nil
	case string:
		var err error
		d.Interval, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// UnmarshalJSON parses the duration from JSON.
func (d *TerminalConfig) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	session_timeout, ok := v["session_timeout"]
	if !ok {
		return errors.New("missing value")
	}
	switch value := session_timeout.(type) {
	case float64:
		d.SessionTimeout = time.Duration(value)
		return nil
	case string:
		var err error
		d.SessionTimeout, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
