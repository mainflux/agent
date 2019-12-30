// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/agent/internal/app/agent/register"
	"github.com/mainflux/agent/internal/pkg/config"
	"github.com/mainflux/agent/pkg/edgex"
	export "github.com/mainflux/export/pkg/config"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/senml"
)

const (
	Path = "./config.toml"
)

var (
	// errInvalidCommand indicates malformed command
	errInvalidCommand = errors.New("invalid command")

	// ErrMalformedEntity indicates malformed entity specification
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrInvalidQueryParams indicates malformed URL
	ErrInvalidQueryParams = errors.New("invalid query params")

	// errUnknownCommand indicates that command is not found
	errUnknownCommand = errors.New("Unknown command")
)

// Service specifies API for publishing messages and subscribing to topics.
type Service interface {
	// Execute command
	Execute(string, string) (string, error)

	// Control command
	Control(string, string) error

	// Update configuration file
	AddConfig(config.Config) error

	// View returns Config struct created from config file
	ViewConfig() config.Config

	// Saves config file for service
	ServiceConfig(string, string) error

	// View returns service list
	ViewServices() map[string]*register.Application

	// Publish message
	Publish(string, string) error
}

var _ Service = (*agent)(nil)

type agent struct {
	mqttClient  paho.Client
	config      *config.Config
	edgexClient edgex.Client
	logger      log.Logger
	register    register.Service
}

// New returns agent service implementation.
func New(mc paho.Client, cfg *config.Config, ec edgex.Client, reg register.Service, logger log.Logger) Service {
	return &agent{
		mqttClient:  mc,
		edgexClient: ec,
		config:      cfg,
		logger:      logger,
		register:    reg,
	}
}

func (a *agent) Execute(uuid, cmd string) (string, error) {
	cmdArr := strings.Split(strings.Replace(cmd, " ", "", -1), ",")
	if len(cmd) < 1 {
		return "", errInvalidCommand
	}

	out, err := exec.Command(cmdArr[0], cmdArr[1:]...).CombinedOutput()
	if err != nil {
		return "", err
	}

	payload, err := encodeSenML(uuid, cmdArr[0], string(out))
	if err != nil {
		return "", err
	}

	if err := a.Publish(a.config.Agent.Channels.Control, string(payload)); err != nil {
		return "", err
	}

	return string(payload), nil
}

func (a *agent) Control(uuid, cmdStr string) error {
	cmdArgs := strings.Split(strings.Replace(cmdStr, " ", "", -1), ",")
	if len(cmdArgs) < 2 {
		return errInvalidCommand
	}

	var resp string
	var err error

	cmd := cmdArgs[0]
	switch cmd {
	case "edgex-operation":
		resp, err = a.edgexClient.PushOperation(cmdArgs[1:])
	case "edgex-config":
		resp, err = a.edgexClient.FetchConfig(cmdArgs[1:])
	case "edgex-metrics":
		resp, err = a.edgexClient.FetchMetrics(cmdArgs[1:])
	case "edgex-ping":
		resp, err = a.edgexClient.Ping()
	default:
		err = errUnknownCommand
	}

	if err != nil {
		return err
	}

	payload, err := encodeSenML(uuid, cmd, resp)
	if err != nil {
		return err
	}

	if err := a.Publish(a.config.Agent.Channels.Control, string(payload)); err != nil {
		return err
	}

	return nil
}

func (a *agent) ServiceConfig(uuid, cmdStr string) error {
	cmdArgs := strings.Split(strings.Replace(cmdStr, " ", "", -1), ",")
	if len(cmdArgs) < 2 {
		return errInvalidCommand
	}

	fileName := cmdArgs[0]
	fileCont := cmdArgs[1]
	c := &export.Config{}

	c.ReadFromB([]byte(fileCont))
	c.File = fileName
	c.Save()

	return nil
}

func (a *agent) AddConfig(c config.Config) error {
	return c.Save()
}

func (a *agent) ViewConfig() config.Config {
	return *a.config
}

func (a *agent) ViewServices() map[string]*register.Application {
	return a.register.Applications()
}

func (a *agent) Publish(crtlChan, payload string) error {
	topic := fmt.Sprintf("channels/%s/messages/res", crtlChan)
	token := a.mqttClient.Publish(topic, 0, false, payload)
	token.Wait()

	return token.Error()
}

func encodeSenML(bn, n, sv string) ([]byte, error) {
	s := senml.Pack{
		Records: []senml.Record{
			senml.Record{
				BaseName:    bn,
				Name:        n,
				StringValue: &sv,
			},
		},
	}

	payload, err := senml.Encode(s, senml.JSON)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
