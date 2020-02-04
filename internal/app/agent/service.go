// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/agent/internal/app/agent/services"
	"github.com/mainflux/agent/internal/pkg/config"
	"github.com/mainflux/agent/pkg/edgex"
	export "github.com/mainflux/export/pkg/config"
	"github.com/mainflux/mainflux/errors"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/senml"
	"github.com/nats-io/go-nats"
)

const (
	Path     = "./config.toml"
	Hearbeat = "heartbeat.*"
	Commands = "commands"
	Config   = "config"

	view = "view"
	save = "save"
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

	// errNatsSubscribing indicates problem with sub to topic for heartbeat
	errNatsSubscribing = errors.New("failed to subscribe to heartbeat topic")

	// errNoSuchService indicates service not supported
	errNoSuchService = errors.New("no such service")
)

// Service specifies API for publishing messages and subscribing to topics.
type Service interface {
	// Execute command
	Execute(string, string) (string, error)

	// Control command
	Control(string, string) error

	// Update configuration file
	AddConfig(config.Config) error

	// Config returns Config struct created from config file
	Config() config.Config

	// Saves config file
	ServiceConfig(uuid, cmdStr string) error

	// Services returns service list
	Services() []services.Service

	// Publish message
	Publish(string, string) error
}

var _ Service = (*agent)(nil)

type agent struct {
	mqttClient  paho.Client
	config      *config.Config
	edgexClient edgex.Client
	logger      log.Logger
	nats        *nats.Conn
	svcs        map[string]*services.Service
}

// New returns agent service implementation.
func New(mc paho.Client, cfg *config.Config, ec edgex.Client, nc *nats.Conn, logger log.Logger) (Service, errors.Error) {
	ag := &agent{
		mqttClient:  mc,
		edgexClient: ec,
		config:      cfg,
		nats:        nc,
		logger:      logger,
		svcs:        make(map[string]*services.Service),
	}

	_, err := ag.nats.Subscribe(Hearbeat, func(msg *nats.Msg) {
		sub := msg.Subject
		tok := strings.Split(sub, ".")
		if len(tok) < 2 {
			ag.logger.Error(fmt.Sprintf("Failed: Subject has incorrect length %s" + sub))
			return
		}
		svcname := tok[1]
		// Service name is extracted from the subtopic
		// if there is multiple instances of the same service
		// we will have to add another distinction
		if _, ok := ag.svcs[svcname]; !ok {
			svc := services.NewService(svcname)
			ag.svcs[svcname] = svc
			ag.logger.Info(fmt.Sprintf("Services '%s' registered", svcname))
		}
		serv := ag.svcs[svcname]
		serv.Update()
	})

	if err != nil {
		return ag, errors.Wrap(errNatsSubscribing, err)
	}

	return ag, nil

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

	return a.processResponse(uuid, cmd, resp)
}

// Message for this command
// [{"bn":"1:", "n":"services", "vs":"view"}]
// [{"bn":"1:", "n":"config", "vs":"save, export, filename, filecontent"}]
// config_file_content is base64 encoded marshaled structure representing service conf
// Example of creation:
// 	b, _ := toml.Marshal(cfg)
// 	config_file_content := base64.StdEncoding.EncodeToString(b)
func (a *agent) ServiceConfig(uuid, cmdStr string) error {
	cmdArgs := strings.Split(strings.Replace(cmdStr, " ", "", -1), ",")
	if len(cmdArgs) < 1 {
		return errInvalidCommand
	}
	resp := ""
	cmd := cmdArgs[0]

	switch cmd {
	case view:
		services, err := json.Marshal(a.Services())
		if err != nil {
			return err
		}
		resp = string(services)
	case save:
		if len(cmdArgs) < 4 {
			return errInvalidCommand
		}
		service := cmdArgs[1]
		fileName := cmdArgs[2]
		fileCont := cmdArgs[3]
		if err := a.saveConfig(service, fileName, fileCont); err != nil {
			return err
		}
	}
	return a.processResponse(uuid, cmd, resp)
}

func (a *agent) processResponse(uuid, cmd, resp string) error {
	payload, err := encodeSenML(uuid, cmd, resp)
	if err != nil {
		return err
	}
	if err := a.Publish(a.config.Agent.Channels.Control, string(payload)); err != nil {
		return err
	}
	return nil
}

func (a *agent) saveConfig(service, fileName, fileCont string) error {
	switch service {
	case "export":
		content, err := base64.StdEncoding.DecodeString(fileCont)
		if err != nil {
			return err
		}

		c := &export.Config{}
		if err := c.ReadBytes([]byte(content)); err != nil {
			return err
		}
		c.File = fileName
		if err := c.Save(); err != nil {
			return err
		}

	default:
		return errNoSuchService
	}

	return a.nats.Publish(fmt.Sprintf("%s.%s.%s", Commands, service, Config), []byte(""))
}

func (a *agent) AddConfig(c config.Config) error {
	return c.Save()
}

func (a *agent) Config() config.Config {
	return *a.config
}

func (a *agent) Services() []services.Service {
	services := [](services.Service){}
	keys := []string{}
	for k := range a.svcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		service := a.svcs[key]
		services = append(services, *service)
	}
	return services
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
