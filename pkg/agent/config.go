// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

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
	ClientCert  string          `json:"client_cert" toml:"client_cert"`
	ClientKey   string          `json:"client_key" toml:"client_key"`
	CaCert      string          `json:"ca_cert" toml:"ca_cert"`
}

type HeartbeatConf struct {
	Interval time.Duration `toml:"interval"`
}

type TerminalConf struct {
	SessionTimeout time.Duration `toml:"session_timeout" json:"session_timeout"`
}

type Config struct {
	Server    ServerConf    `toml:"server" json:"server"`
	Terminal  TerminalConf  `toml:"terminal" json:"terminal"`
	Heartbeat HeartbeatConf `toml:"heartbeat" json:"heartbeat"`
	Channels  ChanConf      `toml:"channels" json:"channels"`
	Edgex     EdgexConf     `toml:"edgex" json:"edgex"`
	Log       LogConf       `toml:"log" json:"log"`
	MQTT      MQTTConf      `toml:"mqtt" json:"mqtt"`
	File      string
}

func NewConfig(sc ServerConf, cc ChanConf, ec EdgexConf, lc LogConf, mc MQTTConf, hc HeartbeatConf, tc TerminalConf, file string) Config {
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

// Save - store config in a file
func SaveConfig(c Config) error {
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
func ReadConfig(file string) (Config, error) {
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

// UnmarshalJSON parses the duration from JSON
func (d *HeartbeatConf) UnmarshalJSON(b []byte) error {
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

// UnmarshalJSON parses the duration from JSON
func (d *TerminalConf) UnmarshalJSON(b []byte) error {
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
