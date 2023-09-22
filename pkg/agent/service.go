// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/agent/pkg/edgex"
	"github.com/mainflux/agent/pkg/encoder"
	"github.com/mainflux/agent/pkg/terminal"

	exp "github.com/mainflux/export/pkg/config"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/pkg/messaging"
)

const (
	Path     = "./config.toml"
	Hearbeat = "channels.heartbeat.>"
	Commands = "commands"
	config   = "config"

	view = "view"
	save = "save"

	char    = "c"
	open    = "open"
	close   = "close"
	control = "control"
	data    = "data"

	export = "export"

	pubSubID = "agent"
)

var (
	// errInvalidCommand indicates malformed command.
	errInvalidCommand = errors.New("invalid command")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrInvalidQueryParams indicates malformed URL.
	ErrInvalidQueryParams = errors.New("invalid query params")

	// errUnknownCommand indicates that command is not found.
	errUnknownCommand = errors.New("Unknown command")

	// errNatsSubscribing indicates problem with sub to topic for heartbeat.
	errNatsSubscribing = errors.New("failed to subscribe to heartbeat topic")

	// errNoSuchService indicates service not supported.
	errNoSuchService = errors.New("no such service")

	// errFailedEncode indicates error in encoding.
	errFailedEncode = errors.New("failed to encode")

	// errFailedToPublish.
	errFailedToPublish = errors.New("failed to publish")

	// errEdgexFailed.
	errEdgexFailed = errors.New("failed to execute edgex operation")

	// errFailedExecute.
	errFailedExecute = errors.New("failed to execute command")

	// errFailedToCreateTerminalSession.
	errFailedToCreateTerminalSession = errors.New("failed to create terminal session")

	// errNoSuchTerminalSession terminal session doesnt exist error on closing.
	errNoSuchTerminalSession = errors.New("no such terminal session")
)

// Service specifies API for publishing messages and subscribing to topics.
type Service interface {
	// Execute command.
	Execute(string, string) (string, error)

	// Control command.
	Control(string, string) error

	// Update configuration file.
	AddConfig(Config) error

	// Config returns Config struct created from config file.
	Config() Config

	// Saves config file.
	ServiceConfig(ctx context.Context, uuid, cmdStr string) error

	// Services returns service list.
	Services() []Info

	// Terminal used for terminal control of gateway.
	Terminal(string, string) error

	// Publish message.
	Publish(string, string) error
}

var _ Service = (*agent)(nil)

type agent struct {
	mqttClient  paho.Client
	config      *Config
	edgexClient edgex.Client
	logger      log.Logger
	broker      messaging.PubSub
	svcs        map[string]Heartbeat
	terminals   map[string]terminal.Session
}

func (ag *agent) handle(ctx context.Context, pub messaging.Publisher, logger log.Logger, cfg HeartbeatConfig) handleFunc {
	return func(msg *messaging.Message) error {
		sub := msg.Channel
		tok := strings.Split(sub, ".")
		if len(tok) < 3 {
			ag.logger.Error(fmt.Sprintf("failed: subject has incorrect length %s", sub))
			return fmt.Errorf("failed: subject has incorrect length %s", sub)
		}
		svcname := tok[1]
		svctype := tok[2]
		// Service name is extracted from the subtopic
		// if there is multiple instances of the same service
		// we will have to add another distinction.
		if _, ok := ag.svcs[svcname]; !ok {
			svc := NewHeartbeat(ctx, svcname, svctype, cfg.Interval)
			ag.svcs[svcname] = svc
			ag.logger.Info(fmt.Sprintf("Services '%s-%s' registered", svcname, svctype))
		}
		serv := ag.svcs[svcname]
		serv.Update()
		return nil
	}
}

type handleFunc func(msg *messaging.Message) error

func (h handleFunc) Handle(msg *messaging.Message) error {
	return h(msg)
}

func (h handleFunc) Cancel() error {
	return nil
}

// New returns agent service implementation.
func New(ctx context.Context, mc paho.Client, cfg *Config, ec edgex.Client, broker messaging.PubSub, logger log.Logger) (Service, error) {
	ag := &agent{
		mqttClient:  mc,
		edgexClient: ec,
		config:      cfg,
		broker:      broker,
		logger:      logger,
		svcs:        make(map[string]Heartbeat),
		terminals:   make(map[string]terminal.Session),
	}

	if cfg.Heartbeat.Interval <= 0 {
		ag.logger.Error(fmt.Sprintf("invalid heartbeat interval %d", cfg.Heartbeat.Interval))
	}

	err := ag.broker.Subscribe(ctx, pubSubID, Hearbeat, ag.handle(ctx, ag.broker, logger, cfg.Heartbeat))

	if err != nil {
		return ag, errors.Wrap(errNatsSubscribing, err)
	}

	return ag, nil

}

func (a *agent) Execute(uuid, cmd string) (string, error) {
	cmdArr := strings.Split(strings.ReplaceAll(cmd, " ", ""), ",")
	if len(cmdArr) < 2 {
		return "", errInvalidCommand
	}

	out, err := exec.Command(cmdArr[0], cmdArr[1:]...).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(errFailedExecute, err)
	}

	payload, err := encoder.EncodeSenML(uuid, cmdArr[0], string(out))
	if err != nil {
		return "", errors.Wrap(errFailedEncode, err)
	}

	if err := a.Publish(control, string(payload)); err != nil {
		return "", errors.Wrap(errFailedToPublish, err)
	}

	return string(payload), nil
}

func (a *agent) Control(uuid, cmdStr string) error {
	cmdArgs := strings.Split(strings.ReplaceAll(cmdStr, " ", ""), ",")
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
		return errors.Wrap(errEdgexFailed, err)
	}

	return a.processResponse(uuid, cmd, resp)
}

// Message for this command
// [{"bn":"1:", "n":"services", "vs":"view"}]
// [{"bn":"1:", "n":"config", "vs":"save, export, filename, filecontent"}]
// config_file_content is base64 encoded marshaled structure representing service conf
// Example of creation:
//
//	b, _ := toml.Marshal(cfg)
//	config_file_content := base64.StdEncoding.EncodeToString(b).
func (a *agent) ServiceConfig(ctx context.Context, uuid, cmdStr string) error {
	cmdArgs := strings.Split(strings.ReplaceAll(cmdStr, " ", ""), ",")
	if len(cmdArgs) < 1 {
		return errInvalidCommand
	}
	resp := ""
	cmd := cmdArgs[0]
	switch cmd {
	case view:
		services, err := json.Marshal(a.Services())
		if err != nil {
			return errors.New(err.Error())
		}
		resp = string(services)
	case save:
		if len(cmdArgs) < 4 {
			return errInvalidCommand
		}
		service := cmdArgs[1]
		fileName := cmdArgs[2]
		fileCont := cmdArgs[3]
		if err := a.saveConfig(ctx, service, fileName, fileCont); err != nil {
			return err
		}
	}
	return a.processResponse(uuid, cmd, resp)
}

func (a *agent) Terminal(uuid, cmdStr string) error {
	b, err := base64.StdEncoding.DecodeString(cmdStr)
	if err != nil {
		return errors.New(err.Error())
	}
	cmdArgs := strings.Split(string(b), ",")
	if len(cmdArgs) < 1 {
		return errInvalidCommand
	}

	cmd := cmdArgs[0]
	ch := ""
	if len(cmdArgs) > 1 {
		ch = cmdArgs[1]
	}
	switch cmd {
	case char:
		if err := a.terminalWrite(uuid, ch); err != nil {
			return err
		}
	case open:
		if err := a.terminalOpen(uuid, a.config.Terminal.SessionTimeout); err != nil {
			return err
		}
	case close:
		if err := a.terminalClose(uuid); err != nil {
			return err
		}
	}
	return nil
}

func (a *agent) terminalOpen(uuid string, timeout time.Duration) error {
	if _, ok := a.terminals[uuid]; !ok {
		term, err := terminal.NewSession(uuid, timeout, a.Publish, a.logger)
		if err != nil {
			return errors.Wrap(errors.Wrap(errFailedToCreateTerminalSession, fmt.Errorf(" for %s", uuid)), err)
		}
		a.terminals[uuid] = term
		go func() {
			for range term.IsDone() {
				// Terminal is inactive, should be closed.
				a.logger.Debug((fmt.Sprintf("Closing terminal session %s", uuid)))
				a.terminalClose(uuid)
				delete(a.terminals, uuid)
				return
			}
		}()
	}
	a.logger.Debug(fmt.Sprintf("Opened terminal session %s", uuid))
	return nil
}

func (a *agent) terminalClose(uuid string) error {
	if _, ok := a.terminals[uuid]; ok {
		delete(a.terminals, uuid)
		a.logger.Debug(fmt.Sprintf("Terminal session: %s closed", uuid))
		return nil
	}
	return errors.Wrap(errNoSuchTerminalSession, fmt.Errorf("session :%s", uuid))
}

func (a *agent) terminalWrite(uuid, cmd string) error {
	if err := a.terminalOpen(uuid, a.config.Terminal.SessionTimeout); err != nil {
		return err
	}
	term := a.terminals[uuid]
	p := []byte(cmd)
	return term.Send(p)
}

func (a *agent) processResponse(uuid, cmd, resp string) error {
	payload, err := encoder.EncodeSenML(uuid, cmd, resp)
	if err != nil {
		return errors.Wrap(errFailedEncode, err)
	}
	if err := a.Publish(control, string(payload)); err != nil {
		return errors.Wrap(errFailedToPublish, err)
	}
	return nil
}

func (a *agent) saveConfig(ctx context.Context, service, fileName, fileCont string) error {
	switch service {
	case export:
		content, err := base64.StdEncoding.DecodeString(fileCont)
		if err != nil {
			return errors.New(err.Error())
		}
		c, err := exp.ReadBytes([]byte(content))
		if err != nil {
			return errors.New(err.Error())
		}
		c.File = fileName
		if err := exp.Save(c); err != nil {
			return errors.New(err.Error())
		}

	default:
		return errNoSuchService
	}

	return a.broker.Publish(ctx, fmt.Sprintf("%s.%s.%s", Commands, service, config), &messaging.Message{})
}

func (a *agent) AddConfig(c Config) error {
	err := SaveConfig(c)
	return errors.New(err.Error())
}

func (a *agent) Config() Config {
	return *a.config
}

func (a *agent) Services() []Info {
	svcInfos := []Info{}
	keys := []string{}
	for k := range a.svcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		service := a.svcs[key].Info()
		svcInfos = append(svcInfos, service)
	}
	return svcInfos
}

func (a *agent) Publish(t, payload string) error {
	topic := a.getTopic(t)
	mqtt := a.config.MQTT
	token := a.mqttClient.Publish(topic, mqtt.QoS, mqtt.Retain, payload)
	token.Wait()
	err := token.Error()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (a *agent) getTopic(topic string) (t string) {
	switch topic {
	case control:
		t = fmt.Sprintf("channels/%s/messages/res", a.config.Channels.Control)
	case data:
		t = fmt.Sprintf("channels/%s/messages/res", a.config.Channels.Data)
	default:
		t = fmt.Sprintf("channels/%s/messages/res/%s", a.config.Channels.Control, topic)
	}
	return t
}
