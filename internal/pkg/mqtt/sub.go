// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"fmt"
	"strings"

	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/senml"

	paho "github.com/eclipse/paho.mqtt.golang"
)

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
}

// NewBroker returns new MQTT broker instance.
func NewBroker(svc agent.Service, client paho.Client, log logger.Logger) MqttBroker {
	return &broker{
		svc:    svc,
		client: client,
		logger: log,
	}
}

// Subscribe subscribes to the MQTT message broker
func (b *broker) Subscribe(topic string) error {
	s := b.client.Subscribe(topic, 0, b.handleMsg)
	if err := s.Error(); s.Wait() && err != nil {
		return err
	}

	return nil
}

// handleMsg triggered when new message is received on MQTT broker
func (b *broker) handleMsg(mc paho.Client, msg paho.Message) {
	sm, err := senml.Decode(msg.Payload(), senml.JSON)
	if err != nil {
		b.logger.Warn(fmt.Sprintf("SenML decode failed: %s", err))
		return
	}

	cmdType := sm.Records[0].Name
	cmdStr := *sm.Records[0].StringValue
	uuid := strings.TrimSuffix(sm.Records[0].BaseName, ":")

	switch cmdType {
	case "control":
		b.logger.Info(fmt.Sprintf("Control command for uuid %s and command string %s", uuid, cmdStr))
		if err := b.svc.Control(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Control operation failed: %s", err))
		}
	case "exec":
		b.logger.Info(fmt.Sprintf("Execute command for uuid %s and command string %s", uuid, cmdStr))
		if _, err := b.svc.Execute(uuid, cmdStr); err != nil {
			b.logger.Warn(fmt.Sprintf("Execute operation failed: %s", err))
		}
	}
}
