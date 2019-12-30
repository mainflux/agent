// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/app/agent/register"
	"github.com/mainflux/agent/internal/pkg/config"
)

var _ agent.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     agent.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc agent.Service, counter metrics.Counter, latency metrics.Histogram) agent.Service {
	return &metricsMiddleware{
		svc:     svc,
		counter: counter,
		latency: latency,
	}
}

func (ms *metricsMiddleware) Execute(uuid, cmdStr string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "execute").Add(1)
		ms.latency.With("method", "execute").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Execute(uuid, cmdStr)
}

func (ms *metricsMiddleware) Control(uuid, cmdStr string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "control").Add(1)
		ms.latency.With("method", "control").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Control(uuid, cmdStr)
}

func (ms *metricsMiddleware) AddConfig(ec config.Config) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_config").Add(1)
		ms.latency.With("method", "add_config").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddConfig(ec)
}

func (ms *metricsMiddleware) ServiceConfig(uuid, cmdStr string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "service_config").Add(1)
		ms.latency.With("method", "service_config").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ServiceConfig(uuid, cmdStr)
}

func (ms *metricsMiddleware) ViewConfig() config.Config {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_config").Add(1)
		ms.latency.With("method", "view_config").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewConfig()
}

func (ms *metricsMiddleware) ViewServices() map[string]*register.Application {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_services").Add(1)
		ms.latency.With("method", "view_services").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewServices()
}

func (ms *metricsMiddleware) Publish(topic, payload string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "execute").Add(1)
		ms.latency.With("method", "execute").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Publish(topic, payload)
}
