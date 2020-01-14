// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"time"

	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/pkg/config"
	log "github.com/mainflux/mainflux/logger"
)

var _ agent.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    agent.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc agent.Service, logger log.Logger) agent.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm loggingMiddleware) Publish(topic string, payload string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method pub for topic %s and payload %s took %s to complete", topic, payload, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Publish(topic, payload)
}

func (lm loggingMiddleware) Execute(uuid, cmd string) (str string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method exec for uuid %s and cmd %s took %s to complete", uuid, cmd, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Execute(uuid, cmd)
}

func (lm loggingMiddleware) Control(uuid, cmd string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method control for uuid %s and cmd %s took %s to complete", uuid, cmd, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Control(uuid, cmd)
}

func (lm loggingMiddleware) AddConfig(c config.Config) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method AddConfig took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.AddConfig(c)
}

func (lm loggingMiddleware) ViewConfig() config.Config {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method ViewConfig took %s to complete", time.Since(begin))
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewConfig()
}

func (lm loggingMiddleware) ViewApplications() map[string]*agent.Application {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method ViewApplications took %s to complete", time.Since(begin))
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewApplications()
}
