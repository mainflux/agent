package conn

import (
	"context"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/agent/pkg/agent"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/stretchr/testify/assert"
)

// Mocks for testing.
type mockService struct{}

func (m *mockService) Config() agent.Config                         { return agent.Config{} }
func (m *mockService) Services() []agent.Info                       { return []agent.Info{} }
func (m *mockService) Publish(string, string) error                 { return nil }
func (m *mockService) AddConfig(agent.Config) error                 { return nil }
func (m *mockService) Control(uuid, command string) error           { return nil }
func (m *mockService) Execute(uuid, command string) (string, error) { return "", nil }
func (m *mockService) ServiceConfig(ctx context.Context, uuid, command string) error {
	return nil
}
func (m *mockService) Close() error                        { return nil }
func (m *mockService) Terminal(uuid, command string) error { return nil }

type mockMQTTClient struct {
	subscribeErr error
	waitErr      error
}

// AddRoute implements mqtt.Client.
func (*mockMQTTClient) AddRoute(topic string, callback mqtt.MessageHandler) {
	panic("unimplemented")
}

// Connect implements mqtt.Client.
func (*mockMQTTClient) Connect() mqtt.Token {
	panic("unimplemented")
}

// Disconnect implements mqtt.Client.
func (*mockMQTTClient) Disconnect(quiesce uint) {
	panic("unimplemented")
}

// IsConnected implements mqtt.Client.
func (*mockMQTTClient) IsConnected() bool {
	panic("unimplemented")
}

// IsConnectionOpen implements mqtt.Client.
func (*mockMQTTClient) IsConnectionOpen() bool {
	panic("unimplemented")
}

// OptionsReader implements mqtt.Client.
func (*mockMQTTClient) OptionsReader() mqtt.ClientOptionsReader {
	panic("unimplemented")
}

// Publish implements mqtt.Client.
func (*mockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	panic("unimplemented")
}

// SubscribeMultiple implements mqtt.Client.
func (*mockMQTTClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	panic("unimplemented")
}

// Unsubscribe implements mqtt.Client.
func (*mockMQTTClient) Unsubscribe(topics ...string) mqtt.Token {
	panic("unimplemented")
}

func (m *mockMQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return &mockToken{err: m.subscribeErr}
}

func (m *mockMQTTClient) Wait() bool {
	return m.waitErr == nil
}

type mockToken struct {
	err error
}

func (m *mockToken) Wait() bool                     { return true }
func (m *mockToken) WaitTimeout(time.Duration) bool { return true }
func (m *mockToken) Error() error                   { return m.err }
func (m *mockToken) Done() <-chan struct{} {
	x := make(chan struct{})
	return x
}

type mockMessageBroker struct {
	publishErr error
}

// Close implements messaging.PubSub.
func (*mockMessageBroker) Close() error {
	panic("unimplemented")
}

// Subscribe implements messaging.PubSub.
func (*mockMessageBroker) Subscribe(ctx context.Context, id string, topic string, handler messaging.MessageHandler) error {
	panic("unimplemented")
}

// Unsubscribe implements messaging.PubSub.
func (*mockMessageBroker) Unsubscribe(ctx context.Context, id string, topic string) error {
	panic("unimplemented")
}

func (m *mockMessageBroker) Publish(ctx context.Context, topic string, msg *messaging.Message) error {
	return m.publishErr
}

func TestBroker_Subscribe(t *testing.T) {
	svc := &mockService{}
	client := &mockMQTTClient{}
	chann := "test"
	messBroker := &mockMessageBroker{}

	broker := NewBroker(svc, client, chann, messBroker, logger.NewMock())

	assert.NotNil(t, broker)

	ctx := context.Background()
	err := broker.Subscribe(ctx)

	assert.NoError(t, err)
}
