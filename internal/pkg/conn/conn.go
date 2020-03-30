// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package conn

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/senml"
	"github.com/nats-io/nats.go"
	"robpike.io/filter"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	reqTopic  = "req"
	servTopic = "services"
	commands  = "commands"

	control = "control"
	exec    = "exec"
	config  = "config"
	service = "service"
	term    = "term"
)

var channelPartRegExp = regexp.MustCompile(`^channels/([\w\-]+)/messages/services(/[^?]*)?(\?.*)?$`)

var _ MqttBroker = (*broker)(nil)

// MqttBroker represents the MQTT broker.
type MqttBroker interface {
	// Subscribes to given topic and receives events.
	Subscribe() error
}

type broker struct {
	svc     agent.Service
	client  mqtt.Client
	logger  logger.Logger
	nats    *nats.Conn
	channel string
}

// NewBroker returns new MQTT broker instance.
func NewBroker(svc agent.Service, client mqtt.Client, chann string, nats *nats.Conn, log logger.Logger) MqttBroker {

	return &broker{
		svc:     svc,
		client:  client,
		logger:  log,
		nats:    nats,
		channel: chann,
	}

}

// Subscribe subscribes to the MQTT message broker
func (b *broker) Subscribe() error {
	topic := fmt.Sprintf("channels/%s/messages/%s", b.channel, reqTopic)
	s := b.client.Subscribe(topic, 0, b.handleMsg)
	if err := s.Error(); s.Wait() && err != nil {
		return err
	}
	topic = fmt.Sprintf("channels/%s/messages/%s/#", b.channel, servTopic)
	if b.nats != nil {
		n := b.client.Subscribe(topic, 0, b.handleNatsMsg)
		if err := n.Error(); n.Wait() && err != nil {
			return err
		}
	}

	return nil
}

// handleNatsMsg triggered when new message is received on MQTT broker
func (b *broker) handleNatsMsg(mc mqtt.Client, msg mqtt.Message) {
	if topic := extractNatsTopic(msg.Topic()); topic != "" {
		b.nats.Publish(topic, msg.Payload())
	}
}

func extractNatsTopic(topic string) string {
	isEmpty := func(s string) bool {
		return (len(s) == 0)
	}
	channelParts := channelPartRegExp.FindStringSubmatch(topic)
	if len(channelParts) < 3 {
		return ""
	}
	filtered := filter.Drop(strings.Split(channelParts[2], "/"), isEmpty).([]string)
	natsTopic := strings.Join(filtered, ".")

	return fmt.Sprintf("%s.%s", commands, natsTopic)
}

// handleMsg triggered when new message is received on MQTT broker
func (b *broker) handleMsg(mc mqtt.Client, msg mqtt.Message) {
	sm, err := senml.Decode(msg.Payload(), senml.JSON)
	if err != nil {
		b.logger.Warn(fmt.Sprintf("SenML decode failed: %s", err))
		return
	}

	if len(sm.Records) == 0 {
		b.logger.Error(fmt.Sprintf("SenML payload empty: `%s`", string(msg.Payload())))
		return
	}
	cmdType := sm.Records[0].Name
	cmdStr := *sm.Records[0].StringValue
	uuid := strings.TrimSuffix(sm.Records[0].BaseName, ":")

	switch cmdType {
	case control:
		b.logger.Info(fmt.Sprintf("Control command for uuid %s and command string %s", uuid, cmdStr))
		if err := b.svc.Control(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Control operation failed: %s", err))
		}
	case exec:
		b.logger.Info(fmt.Sprintf("Execute command for uuid %s and command string %s", uuid, cmdStr))
		if _, err := b.svc.Execute(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Execute operation failed: %s", err))
		}
	case config:
		b.logger.Info(fmt.Sprintf("Config service for uuid %s and command string %s", uuid, cmdStr))
		if err := b.svc.ServiceConfig(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Execute operation failed: %s", err))
		}
	case service:
		b.logger.Info(fmt.Sprintf("Services view for uuid %s and command string %s", uuid, cmdStr))
		if err := b.svc.ServiceConfig(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Services view operation failed: %s", err))
		}
	case term:
		b.logger.Info(fmt.Sprintf("Services view for uuid %s and command string %s", uuid, cmdStr))
		if err := b.svc.Terminal(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Services view operation failed: %s", err))
		}
	case term:
		b.logger.Info(fmt.Sprintf("Services view for uuid %s and command string %s", uuid, cmdStr))
		if err := b.svc.Terminal(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Services view operation failed: %s", err))
		}
	}

}
