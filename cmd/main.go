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
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/agent/pkg/agent"
	"github.com/mainflux/agent/pkg/agent/api"
	"github.com/mainflux/agent/pkg/bootstrap"
	"github.com/mainflux/agent/pkg/conn"
	"github.com/mainflux/agent/pkg/edgex"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/errors"
	"github.com/mainflux/mainflux/logger"
	nats "github.com/nats-io/nats.go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	defHTTPPort                   = "9000"
	defBootstrapURL               = "http://localhost:8202/things/bootstrap"
	defBootstrapID                = ""
	defBootstrapKey               = ""
	defBootstrapRetries           = "5"
	defBootstrapSkipTLS           = "false"
	defBootstrapRetryDelaySeconds = "10"
	defLogLevel                   = "info"
	defEdgexURL                   = "http://localhost:48090/api/v1/"
	defMqttURL                    = "localhost:1883"
	defCtrlChan                   = ""
	defDataChan                   = ""
	defEncryption                 = "false"
	defMqttUsername               = ""
	defMqttPassword               = ""
	defMqttChannel                = ""
	defMqttSkipTLSVer             = "true"
	defMqttMTLS                   = "false"
	defMqttCA                     = "ca.crt"
	defMqttQoS                    = "0"
	defMqttRetain                 = "false"
	defMqttCert                   = "thing.cert"
	defMqttPrivKey                = "thing.key"
	defConfigFile                 = "config.toml"
	defNatsURL                    = nats.DefaultURL
	defHeartbeatInterval          = "10s"
	defTermSessionTimeout         = "60s"
	envConfigFile                 = "MF_AGENT_CONFIG_FILE"
	envLogLevel                   = "MF_AGENT_LOG_LEVEL"
	envEdgexURL                   = "MF_AGENT_EDGEX_URL"
	envMqttURL                    = "MF_AGENT_MQTT_URL"
	envHTTPPort                   = "MF_AGENT_HTTP_PORT"
	envBootstrapURL               = "MF_AGENT_BOOTSTRAP_URL"
	envBootstrapID                = "MF_AGENT_BOOTSTRAP_ID"
	envBootstrapKey               = "MF_AGENT_BOOTSTRAP_KEY"
	envBootstrapRetries           = "MF_AGENT_BOOTSTRAP_RETRIES"
	envBootstrapSkipTLS           = "MF_AGENT_BOOTSTRAP_SKIP_TLS"
	envBootstrapRetryDelaySeconds = "MF_AGENT_BOOTSTRAP_RETRY_DELAY_SECONDS"
	envCtrlChan                   = "MF_AGENT_CONTROL_CHANNEL"
	envDataChan                   = "MF_AGENT_DATA_CHANNEL"
	envEncryption                 = "MF_AGENT_ENCRYPTION"
	envNatsURL                    = "MF_AGENT_NATS_URL"

	envMqttUsername       = "MF_AGENT_MQTT_USERNAME"
	envMqttPassword       = "MF_AGENT_MQTT_PASSWORD"
	envMqttSkipTLSVer     = "MF_AGENT_MQTT_SKIP_TLS"
	envMqttMTLS           = "MF_AGENT_MQTT_MTLS"
	envMqttCA             = "MF_AGENT_MQTT_CA"
	envMqttQoS            = "MF_AGENT_MQTT_QOS"
	envMqttRetain         = "MF_AGENT_MQTT_RETAIN"
	envMqttCert           = "MF_AGENT_MQTT_CLIENT_CERT"
	envMqttPrivKey        = "MF_AGENT_MQTT_CLIENT_PK"
	envHeartbeatInterval  = "MF_AGENT_HEARTBEAT_INTERVAL"
	envTermSessionTimeout = "MF_AGENT_TERMINAL_SESSION_TIMEOUT"
)

var (
	errFailedToSetupMTLS       = errors.New("Failed to set up mtls certs")
	errFetchingBootstrapFailed = errors.New("Fetching bootstrap failed with error")
	errFailedToReadConfig      = errors.New("Failed to read config")
	errFailedToConfigHeartbeat = errors.New("Failed to configure heartbeat")
)

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to load config: %s", err))
	}

	logger, err := logger.New(os.Stdout, cfg.Log.Level)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to create logger: %s", err))
	}

	cfg, err = loadBootConfig(cfg, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load config: %s", err))
	}

	nc, err := nats.Connect(cfg.Server.NatsURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s %s", err, cfg.Server.NatsURL))
		os.Exit(1)
	}
	defer nc.Close()

	mqttClient, err := connectToMQTTBroker(cfg.MQTT, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	edgexClient := edgex.NewClient(cfg.Edgex.URL, logger)

	svc, err := agent.New(mqttClient, &cfg, edgexClient, nc, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in agent service: %s", err))
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
	b := conn.NewBroker(svc, mqttClient, cfg.Channels.Control, nc, logger)
	go b.Subscribe()

	errs := make(chan error, 3)

	go func() {
		p := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info(fmt.Sprintf("Agent service started, exposed port %s", cfg.Server.Port))
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

func loadEnvConfig() (agent.Config, error) {
	sc := agent.ServerConfig{
		NatsURL: mainflux.Env(envNatsURL, defNatsURL),
		Port:    mainflux.Env(envHTTPPort, defHTTPPort),
	}
	cc := agent.ChanConfig{
		Control: mainflux.Env(envCtrlChan, defCtrlChan),
		Data:    mainflux.Env(envDataChan, defDataChan),
	}
	interval, err := time.ParseDuration(mainflux.Env(envHeartbeatInterval, defHeartbeatInterval))
	if err != nil {
		return agent.Config{}, errors.Wrap(errFailedToConfigHeartbeat, err)
	}

	ch := agent.HeartbeatConfig{
		Interval: interval,
	}
	termSessionTimeout, err := time.ParseDuration(mainflux.Env(envTermSessionTimeout, defTermSessionTimeout))
	if err != nil {
		return agent.Config{}, err
	}
	ct := agent.TerminalConfig{
		SessionTimeout: termSessionTimeout,
	}
	ec := agent.EdgexConfig{URL: mainflux.Env(envEdgexURL, defEdgexURL)}
	lc := agent.LogConfig{Level: mainflux.Env(envLogLevel, defLogLevel)}

	mtls, err := strconv.ParseBool(mainflux.Env(envMqttMTLS, defMqttMTLS))
	if err != nil {
		mtls = false
	}

	skipTLSVer, err := strconv.ParseBool(mainflux.Env(defMqttSkipTLSVer, envMqttSkipTLSVer))
	if err != nil {
		skipTLSVer = true
	}

	qos, err := strconv.Atoi(mainflux.Env(envMqttQoS, defMqttQoS))
	if err != nil {
		qos = 0
	}

	retain, err := strconv.ParseBool(mainflux.Env(envMqttRetain, defMqttRetain))
	if err != nil {
		retain = false
	}

	mc := agent.MQTTConfig{
		URL:         mainflux.Env(envMqttURL, defMqttURL),
		Username:    mainflux.Env(envMqttUsername, defMqttUsername),
		Password:    mainflux.Env(envMqttPassword, defMqttPassword),
		MTLS:        mtls,
		CAPath:      mainflux.Env(envMqttCA, defMqttCA),
		CertPath:    mainflux.Env(envMqttCert, defMqttCert),
		PrivKeyPath: mainflux.Env(envMqttPrivKey, defMqttPrivKey),
		SkipTLSVer:  skipTLSVer,
		QoS:         byte(qos),
		Retain:      retain,
	}

	file := mainflux.Env(envConfigFile, defConfigFile)
	c := agent.NewConfig(sc, cc, ec, lc, mc, ch, ct, file)
	mc, err = loadCertificate(c.MQTT)
	if err != nil {
		return c, errors.Wrap(errFailedToSetupMTLS, err)
	}

	c.MQTT = mc
	agent.SaveConfig(c)
	return c, nil
}

func loadBootConfig(c agent.Config, logger logger.Logger) (bsc agent.Config, err error) {
	file := mainflux.Env(envConfigFile, defConfigFile)
	skipTLS, err := strconv.ParseBool(mainflux.Env(envBootstrapSkipTLS, defBootstrapSkipTLS))
	bsConfig := bootstrap.Config{
		URL:           mainflux.Env(envBootstrapURL, defBootstrapURL),
		ID:            mainflux.Env(envBootstrapID, defBootstrapID),
		Key:           mainflux.Env(envBootstrapKey, defBootstrapKey),
		Retries:       mainflux.Env(envBootstrapRetries, defBootstrapRetries),
		RetryDelaySec: mainflux.Env(envBootstrapRetryDelaySeconds, defBootstrapRetryDelaySeconds),
		Encrypt:       mainflux.Env(envEncryption, defEncryption),
		SkipTLS:       skipTLS,
	}

	if err := bootstrap.Bootstrap(bsConfig, logger, file); err != nil {
		return c, errors.Wrap(errFetchingBootstrapFailed, err)
	}

	if bsc, err = agent.ReadConfig(file); err != nil {
		return c, errors.Wrap(errFailedToReadConfig, err)
	}

	mc, err := loadCertificate(bsc.MQTT)
	if err != nil {
		return bsc, errors.Wrap(errFailedToSetupMTLS, err)
	}

	if bsc.Heartbeat.Interval <= 0 {
		bsc.Heartbeat.Interval = c.Heartbeat.Interval
	}

	if bsc.Terminal.SessionTimeout <= 0 {
		bsc.Terminal.SessionTimeout = c.Terminal.SessionTimeout
	}

	bsc.MQTT = mc
	return bsc, nil
}

func connectToMQTTBroker(conf agent.MQTTConfig, logger logger.Logger) (mqtt.Client, error) {
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

func loadCertificate(cnfg agent.MQTTConfig) (c agent.MQTTConfig, err error) {
	var caByte []byte
	var cc []byte
	var pk []byte
	c = cnfg

	cert := tls.Certificate{}
	if !c.MTLS {
		return c, nil
	}
	// Load CA cert from file
	if c.CAPath != "" {
		caFile, err := os.Open(c.CAPath)
		defer caFile.Close()
		if err != nil {
			return c, err
		}
		caByte, err = ioutil.ReadAll(caFile)
		if err != nil {
			return c, err
		}
	}
	// Load CA cert from string if file not present
	if len(caByte) == 0 && c.CaCert != "" {
		caByte, err = ioutil.ReadAll(strings.NewReader(c.CaCert))
		if err != nil {
			return c, err
		}
	}
	// Load client certificate from file if present
	if c.CertPath != "" {
		clientCert, err := os.Open(c.CertPath)
		defer clientCert.Close()
		if err != nil {
			return c, err
		}
		cc, err = ioutil.ReadAll(clientCert)
		if err != nil {
			return c, err
		}
	}
	// Load client certificate from string if file not present
	if len(cc) == 0 && c.ClientCert != "" {
		cc, err = ioutil.ReadAll(strings.NewReader(c.ClientCert))
		if err != nil {
			return c, err
		}
	}
	// Load private key of client certificate from file
	if c.PrivKeyPath != "" {
		privKey, err := os.Open(c.PrivKeyPath)
		defer privKey.Close()
		if err != nil {
			return c, err
		}
		pk, err = ioutil.ReadAll((privKey))
		if err != nil {
			return c, err
		}
	}
	// Load private key of client certificate from string
	if len(pk) == 0 && c.ClientKey != "" {
		pk, err = ioutil.ReadAll(strings.NewReader(c.ClientKey))
		if err != nil {
			return c, err
		}
	}

	cert, err = tls.X509KeyPair([]byte(c.ClientCert), []byte(c.ClientKey))
	if err != nil {
		return c, err
	}
	c.Cert = cert
	c.CA = caByte
	return c, nil
}
