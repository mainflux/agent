package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	paho "github.com/eclipse/paho.mqtt.golang"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/app/agent/api"
	"github.com/mainflux/agent/internal/pkg/bootstrap"
	"github.com/mainflux/agent/internal/pkg/config"
	"github.com/mainflux/agent/internal/pkg/mqtt"
	"github.com/mainflux/agent/pkg/edgex"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	nats "github.com/nats-io/go-nats"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	defHTTPPort                   = "9000"
	defBootstrapURL               = "http://localhost:8202/things/bootstrap"
	defBootstrapID                = "75-7E-36-73-54-EC"
	defBootstrapKey               = "2cf9cb98-5ae9-42d8-bb21-1b6af97a490c"
	defBootstrapRetries           = "5"
	defBootstrapRetryDelaySeconds = "10"
	defLogLevel                   = "info"
	defEdgexURL                   = "http://localhost:48090/api/v1/"
	defMqttURL                    = "localhost:1883"
	defThingID                    = "2dce1d65-73b4-4020-bfe3-403d851386e7"
	defThingKey                   = "1ff0d0f0-ea04-4fbb-83c4-c10b110bf566"
	defCtrlChan                   = "f36c3733-95a3-481c-a314-4125e03d8993"
	defDataChan                   = "ea353dac-0298-4fbb-9e5d-501e3699949c"
	defEncryption                 = "false"
	defConfigFile                 = "config.toml"
	defNatsURL                    = nats.DefaultURL

	envConfigFile                 = "MF_AGENT_CONFIG_FILE"
	envLogLevel                   = "MF_AGENT_LOG_LEVEL"
	envEdgexURL                   = "MF_AGENT_EDGEX_URL"
	envMqttURL                    = "MF_AGENT_MQTT_URL"
	envHTTPPort                   = "MF_AGENT_HTTP_PORT"
	envBootstrapURL               = "MF_AGENT_BOOTSTRAP_URL"
	envBootstrapID                = "MF_AGENT_BOOTSTRAP_ID"
	envBootstrapKey               = "MF_AGENT_BOOTSTRAP_KEY"
	envBootstrapRetries           = "MF_AGENT_BOOTSTRAP_RETRIES"
	envBootstrapRetryDelaySeconds = "MF_AGENT_BOOTSTRAP_RETRY_DELAY_SECONDS"
	envThingID                    = "MF_AGENT_THING_ID"
	envThingKey                   = "MF_AGENT_THING_KEY"
	envCtrlChan                   = "MF_AGENT_CONTROL_CHANNEL"
	envDataChan                   = "MF_AGENT_DATA_CHANNEL"
	envEncryption                 = "MF_AGENT_ENCRYPTION"
	envNatsURL                    = "MF_AGENT_NATS_URL"
)

func main() {
	logger, err := logger.New(os.Stdout, defLogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to create logger: %s", err.Error()))
	}

	cfg, err := loadConfig(logger)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to load config: %s", err.Error()))
	}

	nc, err := nats.Connect(cfg.Agent.Server.NatsURL)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to NATS: %s %s", err, cfg.Agent.Server.NatsURL))
	}
	defer nc.Close()

	mqttClient := connectToMQTTBroker(cfg.Agent.MQTT.URL, cfg.Agent.Thing.ID, cfg.Agent.Thing.Key, logger)
	edgexClient := edgex.NewClient(cfg.Agent.Edgex.URL, logger)

	svc, err := agent.New(mqttClient, &cfg, edgexClient, nc, logger)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error in agent service: %s", err.Error()))
	}

	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "agent",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "agent",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	go subscribeToMQTTBroker(svc, mqttClient, cfg.Agent.Channels.Control, nc, logger)

	errs := make(chan error, 3)

	go func() {
		p := fmt.Sprintf(":%s", cfg.Agent.Server.Port)
		logger.Info(fmt.Sprintf("Agent service started, exposed port %s", cfg.Agent.Server.Port))
		errs <- http.ListenAndServe(p, api.MakeHandler(svc))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Agent terminated: %s", err))
}

func loadConfig(logger logger.Logger) (config.Config, error) {
	file := mainflux.Env(envConfigFile, defConfigFile)

	bcfg := bootstrap.Config{
		URL:           mainflux.Env(envBootstrapURL, defBootstrapURL),
		ID:            mainflux.Env(envBootstrapID, defBootstrapID),
		Key:           mainflux.Env(envBootstrapKey, defBootstrapKey),
		Retries:       mainflux.Env(envBootstrapRetries, defBootstrapRetries),
		RetryDelaySec: mainflux.Env(envBootstrapRetryDelaySeconds, defBootstrapRetryDelaySeconds),
		Encrypt:       mainflux.Env(envEncryption, defEncryption),
	}
	if err := bootstrap.Bootstrap(bcfg, logger, file); err != nil {
		logger.Error(fmt.Sprintf("Fetching bootstrap failed with error: %s", err))
		return config.Config{}, err
	}

	sc := config.ServerConf{
		NatsURL: mainflux.Env(envNatsURL, defNatsURL),
		Port:    mainflux.Env(envLogLevel, defLogLevel),
	}
	tc := config.ThingConf{
		ID:  mainflux.Env(envThingID, defThingID),
		Key: mainflux.Env(envThingKey, defThingKey),
	}
	cc := config.ChanConf{
		Control: mainflux.Env(envCtrlChan, defCtrlChan),
		Data:    mainflux.Env(envDataChan, defDataChan),
	}
	ec := config.EdgexConf{URL: mainflux.Env(envEdgexURL, defEdgexURL)}
	lc := config.LogConf{Level: mainflux.Env(envLogLevel, defLogLevel)}
	mc := config.MQTTConf{URL: mainflux.Env(envMqttURL, defMqttURL)}

	c := config.New(sc, tc, cc, ec, lc, mc, file)

	if err := c.Read(); err != nil {
		logger.Error(fmt.Sprintf("Failed to read config:  %s", err))
		return config.Config{}, err
	}

	return *c, nil
}

func connectToMQTTBroker(mqttURL, thingID, thingKey string, logger logger.Logger) paho.Client {
	opts := paho.NewClientOptions()
	opts.AddBroker(mqttURL)
	opts.SetClientID("agent")
	opts.SetUsername(thingID)
	opts.SetPassword(thingKey)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(c paho.Client) {
		logger.Info("Connected to MQTT broker")
	})
	opts.SetConnectionLostHandler(func(c paho.Client, err error) {
		logger.Error(fmt.Sprintf("MQTT connection lost: %s", err.Error()))
		os.Exit(1)
	})

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(fmt.Sprintf("Failed to connect to MQTT broker: %s", token.Error()))
		os.Exit(1)
	}

	return client
}

func subscribeToMQTTBroker(svc agent.Service, mc paho.Client, ctrlChan string, nc *nats.Conn, logger logger.Logger) {
	broker := mqtt.NewBroker(svc, mc, nc, logger)
	topic := fmt.Sprintf("channels/%s/messages", ctrlChan)
	if err := broker.Subscribe(topic); err != nil {
		logger.Error(fmt.Sprintf("Failed to subscribe to MQTT broker: %s", err.Error()))
		os.Exit(1)
	}
	logger.Info("Subscribed to MQTT broker")
}
