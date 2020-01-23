// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/pkg/config"
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

		sc := config.ServerConf{Port: req.agent.server.port}
		cc := config.ChanConf{
			Control: req.agent.channels.control,
			Data:    req.agent.channels.data,
		}
		ec := config.EdgexConf{URL: req.agent.edgex.url}
		lc := config.LogConf{Level: req.agent.log.level}
		mc := config.MQTTConf{
			URL:      req.agent.mqtt.url,
			Username: req.agent.mqtt.username,
			Password: req.agent.mqtt.password,
		}
		a := config.AgentConf{
			Server:   sc,
			Channels: cc,
			Edgex:    ec,
			Log:      lc,
			MQTT:     mc,
		}
		c := config.Config{Agent: a}
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

		res := configRes{agent: c.Agent}

		return res, nil
	}
}

func viewServicesEndpoint(svc agent.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return svc.Services(), nil
	}
}
