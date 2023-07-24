// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	dockertest "github.com/ory/dockertest/v3"
)

const (
	username      = "mainflux-mqtt"
	broker        = "eclipse-mosquitto"
	brokerVersion = "1.6.13"
	poolMaxWait   = 120 * time.Second
)

var (
	natsAddress string
	mqttAddress string
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	container, err := pool.Run("nats", "1.3.0", []string{})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, container)

	address := fmt.Sprintf("%s:%s", "localhost", container.GetPort("4222/tcp"))
	if err := pool.Retry(func() error {
		natsAddress = address
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	mqttContainer, err := pool.Run(broker, brokerVersion, []string{})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	handleInterrupt(pool, mqttContainer)

	address2 := fmt.Sprintf("%s:%s", "localhost", mqttContainer.GetPort("1883/tcp"))
	pool.MaxWait = poolMaxWait

	if err := pool.Retry(func() error {
		mqttAddress = address2
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}
	if err := pool.Purge(mqttContainer); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}

func handleInterrupt(pool *dockertest.Pool, container *dockertest.Resource) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if err := pool.Purge(container); err != nil {
			log.Fatalf("Could not purge container: %s", err)
		}
		os.Exit(0)
	}()
}
