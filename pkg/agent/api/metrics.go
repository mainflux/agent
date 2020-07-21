// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/agent/pkg/agent"
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

func (ms *metricsMiddleware) AddConfig(ec agent.Config) error {
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

func (ms *metricsMiddleware) Config() agent.Config {
	defer func(begin time.Time) {
		ms.counter.With("method", "config").Add(1)
		ms.latency.With("method", "config").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Config()
}

func (ms *metricsMiddleware) Services() []agent.Info {
	defer func(begin time.Time) {
		ms.counter.With("method", "services").Add(1)
		ms.latency.With("method", "services").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Services()
}

func (ms *metricsMiddleware) Publish(topic, payload string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "publish").Add(1)
		ms.latency.With("method", "publish").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Publish(topic, payload)
}

func (ms *metricsMiddleware) Terminal(topic, payload string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "publish").Add(1)
		ms.latency.With("method", "publish").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Terminal(topic, payload)
}
