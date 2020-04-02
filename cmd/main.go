package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/app/agent/api"
	"github.com/mainflux/agent/internal/pkg/bootstrap"
	"github.com/mainflux/agent/internal/pkg/config"
	"github.com/mainflux/agent/internal/pkg/conn"
	"github.com/mainflux/agent/pkg/edgex"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	nats "github.com/nats-io/nats.go"
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
	defCtrlChan                   = "f36c3733-95a3-481c-a314-4125e03d8993"
	defDataChan                   = "ea353dac-0298-4fbb-9e5d-501e3699949c"
	defEncryption                 = "false"
	defMqttUsername               = ""
	defMqttPassword               = ""
	defMqttChannel                = ""
	defMqttSkipTLSVer             = "true"
	defMqttMTLS                   = "false"
	defMqttCA                     = "ca.crt"
	defMqttQoS                    = "0"
	defMqttRetain                 = false
	defMqttCert                   = "thing.cert"
	defMqttPrivKey                = "thing.key"
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

	envMqttUsername   = "MF_AGENT_MQTT_USERNAME"
	envMqttPassword   = "MF_AGENT_MQTT_PASSWORD"
	envMqttSkipTLSVer = "MF_AGENT_MQTT_SKIP_TLS"
	envMqttMTLS       = "MF_AGENT_MQTT_MTLS"
	envMqttCA         = "MF_AGENT_MQTT_CA"
	envMqttQoS        = "MF_AGENT_MQTT_QOS"
	envMqttRetain     = "MF_AGENT_MQTT_RETAIN"
	envMqttCert       = "MF_AGENT_MQTT_CLIENT_CERT"
	envMqttPrivKey    = "MF_AGENT_MQTT_CLIENT_PK"
)

func main() {
	cfg, logger, err := loadConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load config: %s", err.Error()))
		os.Exit(1)
	}

	nc, err := nats.Connect(cfg.Agent.Server.NatsURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s %s", err, cfg.Agent.Server.NatsURL))
		os.Exit(1)
	}
	defer nc.Close()

	mqttClient, err := connectToMQTTBroker(cfg.Agent.MQTT, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	edgexClient := edgex.NewClient(cfg.Agent.Edgex.URL, logger)

	svc, err := agent.New(mqttClient, &cfg, edgexClient, nc, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in agent service: %s", err.Error()))
		os.Exit(1)
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
	b := conn.NewBroker(svc, mqttClient, cfg.Agent.Channels.Control, nc, logger)
	go b.Subscribe()

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

func loadConfig() (config.Config, logger.Logger, error) {
	lc := config.LogConf{Level: mainflux.Env(envLogLevel, defLogLevel)}
	logg, err := logger.New(os.Stdout, lc.Level)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to create logger: %s", err.Error()))
	}
	file := mainflux.Env(envConfigFile, defConfigFile)
	bcfg := bootstrap.Config{
		URL:           mainflux.Env(envBootstrapURL, defBootstrapURL),
		ID:            mainflux.Env(envBootstrapID, defBootstrapID),
		Key:           mainflux.Env(envBootstrapKey, defBootstrapKey),
		Retries:       mainflux.Env(envBootstrapRetries, defBootstrapRetries),
		RetryDelaySec: mainflux.Env(envBootstrapRetryDelaySeconds, defBootstrapRetryDelaySeconds),
		Encrypt:       mainflux.Env(envEncryption, defEncryption),
	}
	if err := bootstrap.Bootstrap(bcfg, logg, file); err != nil {
		logg.Error(fmt.Sprintf("Fetching bootstrap failed with error: %s", err))
		return config.Config{}, logg, err
	}

	sc := config.ServerConf{
		NatsURL: mainflux.Env(envNatsURL, defNatsURL),
		Port:    mainflux.Env(envHTTPPort, defHTTPPort),
	}
	cc := config.ChanConf{
		Control: mainflux.Env(envCtrlChan, defCtrlChan),
		Data:    mainflux.Env(envDataChan, defDataChan),
	}
	ec := config.EdgexConf{URL: mainflux.Env(envEdgexURL, defEdgexURL)}

	mc := config.MQTTConf{
		URL:      mainflux.Env(envMqttURL, defMqttURL),
		Username: mainflux.Env(envMqttUsername, defMqttUsername),
		Password: mainflux.Env(envMqttPassword, defMqttPassword),
	}

	c := config.New(sc, cc, ec, lc, mc, file)
	if err := c.Read(); err != nil {
		logg.Error(fmt.Sprintf("Failed to read config:  %s", err))
		return config.Config{}, logg, err
	}
	logg, err = logger.New(os.Stdout, c.Agent.Log.Level)
	if err != nil {
		log.Fatalf("Failed to recreate logger: %s", err)
	}

	mc, err = loadCertificate(c.Agent.MQTT)
	if err != nil {
		logg.Error(fmt.Sprintf("Failed to set up mtls certs %s", err))
	}
	c.Agent.MQTT = mc

	return *c, logg, nil
}

func connectToMQTTBroker(conf config.MQTTConf, logger logger.Logger) (mqtt.Client, error) {
	name := fmt.Sprintf("agent-%s", conf.Username)
	conn := func(client mqtt.Client) {
		logger.Info(fmt.Sprintf("Client %s connected", name))
	}

	lost := func(client mqtt.Client, err error) {
		logger.Info(fmt.Sprintf("Client %s disconnected", name))
	}

	opts := mqtt.NewClientOptions().
		AddBroker(conf.URL).
		SetClientID(name).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(conn).
		SetConnectionLostHandler(lost)

	if conf.Username != "" && conf.Password != "" {
		opts.SetUsername(conf.Username)
		opts.SetPassword(conf.Password)
	}

	if conf.MTLS {
		cfg := &tls.Config{
			InsecureSkipVerify: conf.SkipTLSVer,
		}

		if conf.CA != nil {
			cfg.RootCAs = x509.NewCertPool()
			cfg.RootCAs.AppendCertsFromPEM(conf.CA)
		}
		if conf.Cert.Certificate != nil {
			cfg.Certificates = []tls.Certificate{conf.Cert}
		}

		cfg.BuildNameToCertificate()
		opts.SetTLSConfig(cfg)
		opts.SetProtocolVersion(4)
	}
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func loadCertificate(cfg config.MQTTConf) (config.MQTTConf, error) {
	c := cfg
	caByte := []byte{}
	cert := tls.Certificate{}
	if !cfg.MTLS {
		return c, nil
	}
	caFile, err := os.Open(cfg.CAPath)
	defer caFile.Close()
	if err != nil {
		return c, err
	}
	caByte, err = ioutil.ReadAll(caFile)
	if err != nil {
		return c, err
	}
	clientCert, err := os.Open(cfg.CertPath)
	defer clientCert.Close()
	if err != nil {
		return c, err
	}
	cc, _ := ioutil.ReadAll(clientCert)
	privKey, err := os.Open(cfg.PrivKeyPath)
	defer clientCert.Close()
	if err != nil {
		return c, err
	}
	pk, _ := ioutil.ReadAll((privKey))
	cert, err = tls.X509KeyPair([]byte(cc), []byte(pk))
	if err != nil {
		return c, err
	}
	cfg.Cert = cert
	cfg.CA = caByte
	return c, nil
}
