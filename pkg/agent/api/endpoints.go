// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/agent/pkg/agent"
)

func pubEndpoint(svc agent.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(pubReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		topic := req.Topic
		payload := req.Payload

		if err := svc.Publish(topic, payload); err != nil {
			return genericRes{}, nil
		}

		return genericRes{
			Service:  "agent",
			Response: "config",
		}, nil
	}
}

func execEndpoint(svc agent.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(execReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		uuid := strings.TrimSuffix(req.BaseName, ":")
		out, err := svc.Execute(uuid, req.Value)
		if err != nil {
			return execRes{}, nil
		}

		resp := execRes{
			BaseName: req.BaseName,
			Name:     "exec",
			Value:    out,
		}
		return resp, nil
	}
}

func addConfigEndpoint(svc agent.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addConfigReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		sc := agent.ServerConfig{Port: req.agent.server.port}
		cc := agent.ChanConfig{
			Control: req.agent.channels.control,
			Data:    req.agent.channels.data,
		}
		ec := agent.EdgexConfig{URL: req.agent.edgex.url}
		lc := agent.LogConfig{Level: req.agent.log.level}
		mc := agent.MQTTConfig{
			URL:      req.agent.mqtt.url,
			Username: req.agent.mqtt.username,
			Password: req.agent.mqtt.password,
		}
		c := agent.Config{
			Server:   sc,
			Channels: cc,
			Edgex:    ec,
			Log:      lc,
			MQTT:     mc,
		}

		if err := svc.AddConfig(c); err != nil {
			return genericRes{}, nil
		}

		return genericRes{
			Service:  "agent",
			Response: "config",
		}, nil
	}
}

func viewConfigEndpoint(svc agent.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		c := svc.Config()
		return c, nil
	}
}

func viewServicesEndpoint(svc agent.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return svc.Services(), nil
	}
}
