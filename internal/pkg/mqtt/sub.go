// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/senml"
	"github.com/nats-io/go-nats"
	"robpike.io/filter"

	paho "github.com/eclipse/paho.mqtt.golang"
)

<<<<<<< HEAD
type cmdType string

const (
	reqTopic  = "req"
	servTopic = "services"
	commands  = "commands"

	control cmdType = "control"
	exec    cmdType = "exec"
	config  cmdType = "config"
)

var channelPartRegExp = regexp.MustCompile(`^channels/([\w\-]+)/messages/services(/[^?]*)?(\?.*)?$`)

=======
const (
	REQ_TOPIC  = "req"
	NATS_TOPIC = "nats"
)

>>>>>>> add nats and support forwarding messages from mqtt to local nats
var _ MqttBroker = (*broker)(nil)

// MqttBroker represents the MQTT broker.
type MqttBroker interface {
	// Subscribes to given topic and receives events.
	Subscribe(string) error
}

type broker struct {
	svc    agent.Service
	client paho.Client
	logger logger.Logger
	nats   *nats.Conn
}

// NewBroker returns new MQTT broker instance.
func NewBroker(svc agent.Service, client paho.Client, nats *nats.Conn, log logger.Logger) MqttBroker {
	return &broker{
		svc:    svc,
		client: client,
		logger: log,
		nats:   nats,
	}
}

// Subscribe subscribes to the MQTT message broker
func (b *broker) Subscribe(topic string) error {
	s := b.client.Subscribe(fmt.Sprintf("%s/%s", topic, reqTopic), 0, b.handleMsg)
	if err := s.Error(); s.Wait() && err != nil {
		return err
	}

	if b.nats != nil {
		n := b.client.Subscribe(fmt.Sprintf("%s/%s/#", topic, servTopic), 0, b.handleNatsMsg)
		if err := n.Error(); n.Wait() && err != nil {
			return err
		}
	}

	return nil
}

// handleNatsMsg triggered when new message is received on MQTT broker
func (b *broker) handleNatsMsg(mc paho.Client, msg paho.Message) {
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
func (b *broker) handleMsg(mc paho.Client, msg paho.Message) {
	sm, err := senml.Decode(msg.Payload(), senml.JSON)
	if err != nil {
		b.logger.Warn(fmt.Sprintf("SenML decode failed: %s", err))
		return
	}

	cmdType := cmdType(sm.Records[0].Name)
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
	}
}
