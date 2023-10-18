// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/agent/pkg/agent"
	"github.com/mainflux/agent/pkg/agent/api"
	"github.com/mainflux/agent/pkg/agent/mocks"

	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging/brokers"
	"github.com/stretchr/testify/assert"
)

type testRequest struct {
	client *http.Client
	method string
	url    string
	body   io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}

	return tr.client.Do(req)
}

func newService(ctx context.Context) (agent.Service, error) {
	opts := paho.NewClientOptions().
		SetUsername(username).
		AddBroker(mqttAddress).
		SetClientID("testing")

	mqttClient := paho.NewClient(opts)
	token := mqttClient.Connect()
	if token.Error() != nil {
		return nil, token.Error()
	}

	edgexClient := mocks.NewEdgexClient()
	config := agent.Config{}
	config.Heartbeat.Interval = time.Second

	logger, err := logger.New(os.Stdout, "debug")
	if err != nil {
		return nil, err
	}

	pubsub, err := brokers.NewPubSub(brokerAddress, "", logger)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Broker: %s %s", err, brokerAddress)
	}
	defer pubsub.Close()

	agentSvc, err := agent.New(ctx, mqttClient, &config, edgexClient, pubsub, logger)
	if err != nil {
		return nil, err
	}

	return agentSvc, nil
}

func newServer(svc agent.Service) *httptest.Server {
	mux := api.MakeHandler(svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestPublish(t *testing.T) {
	svc, err := newService(context.TODO())
	if err != nil {
		t.Errorf("failed to create service: %v", err)
		return
	}
	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()
	data := toJSON(struct {
		Payload string
		Topic   string
	}{
		"payload",
		"topic",
	})

	cases := []struct {
		desc   string
		req    string
		status int
	}{
		{"publish data", data, http.StatusOK},
		{"publish data with invalid data", "}", http.StatusInternalServerError},
	}
	t.Cleanup(func() {
		assert.Nil(t, svc.Close())
	})

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodPost,
			url:    fmt.Sprintf("%s/pub", ts.URL),
			body:   strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}
