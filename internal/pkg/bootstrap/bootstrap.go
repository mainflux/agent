// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"os"

	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/pkg/config"
	export "github.com/mainflux/export/pkg/config"
	errors "github.com/mainflux/mainflux/errors"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/things"
)

const exportConfigFile = "/configs/export/config.toml"

// Config represents the parameters for boostraping
type Config struct {
	URL           string
	ID            string
	Key           string
	Retries       string
	RetryDelaySec string
	Encrypt       string
}

type deviceConfig struct {
	MainfluxID       string           `json:"mainflux_id"`
	MainfluxKey      string           `json:"mainflux_key"`
	MainfluxChannels []things.Channel `json:"mainflux_channels"`
	Content          string           `json:"content"`
}

type infraConfig struct {
	LogLevel     string        `json:"log_level"`
	HTTPPort     string        `json:"http_port"`
	MqttURL      string        `json:"mqtt_url"`
	EdgexURL     string        `json:"edgex_url"`
	NatsURL      string        `json:"nats_url"`
	ExportConfig export.Config `json:"export_config"`
}

// Bootstrap - Retrieve device config
func Bootstrap(cfg Config, logger log.Logger, file string) error {
	retries, err := strconv.ParseUint(cfg.Retries, 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid BOOTSTRAP_RETRIES value: %s", err))
	}

	retryDelaySec, err := strconv.ParseUint(cfg.RetryDelaySec, 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid BOOTSTRAP_RETRY_DELAY_SECONDS value: %s", err))
	}

	logger.Info(fmt.Sprintf("Requesting config for %s from %s", cfg.ID, cfg.URL))

	dc := deviceConfig{}
	for i := 0; i < int(retries); i++ {
		dc, err = getConfig(cfg.ID, cfg.Key, cfg.URL, logger)
		if err == nil {
			break
		}
		logger.Error(fmt.Sprintf("Fetching bootstrap failed with error: %s", err))
		logger.Debug(fmt.Sprintf("Retries remaining: %d. Retrying in %d seconds", retries, retryDelaySec))
		time.Sleep(time.Duration(retryDelaySec) * time.Second)
		if i == int(retries)-1 {
			logger.Warn("Retries exhausted")
			logger.Info(fmt.Sprintf("Continuing with local config"))
			return nil
		}
	}

	logger.Info(fmt.Sprintf("Getting config for %s from %s succeeded",
		cfg.ID, cfg.URL))
	fmt.Printf("cont:%s", dc.Content)
	ic := infraConfig{}
	if err := json.Unmarshal([]byte(dc.Content), &ic); err != nil {
		return errors.New(err.Error())
	}

	saveExportConfig(ic.ExportConfig, logger)

	if len(dc.MainfluxChannels) < 2 {
		return agent.ErrMalformedEntity
	}

	ctrlChan := dc.MainfluxChannels[0].ID
	dataChan := dc.MainfluxChannels[1].ID
	if dc.MainfluxChannels[0].Metadata["type"] == "data" {
		ctrlChan = dc.MainfluxChannels[1].ID
		dataChan = dc.MainfluxChannels[0].ID
	}

	sc := config.ServerConf{
		Port:    ic.HTTPPort,
		NatsURL: ic.NatsURL,
	}

	cc := config.ChanConf{
		Control: ctrlChan,
		Data:    dataChan,
	}
	ec := config.EdgexConf{URL: ic.EdgexURL}
	lc := config.LogConf{Level: ic.LogLevel}
	mc := config.MQTTConf{
		URL:      ic.MqttURL,
		Password: dc.MainfluxKey,
		Username: dc.MainfluxID,
	}

	c := config.New(sc, cc, ec, lc, mc, file)

	return config.Save(c)
}

func saveExportConfig(econf export.Config, logger log.Logger) {
	if econf.File == "" {
		econf.File = exportConfigFile
	}
	exConfFileExist := false
	if _, err := os.Stat(econf.File); err == nil {
		exConfFileExist = true
		logger.Info(fmt.Sprintf("Export config file %s exists", econf.File))
	}
	if !exConfFileExist {
		logger.Info(fmt.Sprintf("Saving export config file %s", econf.File))
		if err := export.Save(econf); err != nil {
			logger.Error(fmt.Sprintf("Failed to save export config file %s", err))
		}
	}
}

func getConfig(bsID, bsKey, bsSvrURL string, logger log.Logger) (deviceConfig, error) {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		logger.Error(err.Error())
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("%s/%s", bsSvrURL, bsID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return deviceConfig{}, errors.New(err.Error())
	}

	req.Header.Add("Authorization", bsKey)
	resp, err := client.Do(req)
	if err != nil {
		return deviceConfig{}, errors.New(err.Error())
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return deviceConfig{}, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return deviceConfig{}, errors.New(err.Error())
	}
	defer resp.Body.Close()
	fmt.Printf("url:%s", url)
	fmt.Printf("body:%s", string(body))
	dc := deviceConfig{}
	if err := json.Unmarshal(body, &dc); err != nil {
		return deviceConfig{}, errors.New(err.Error())
	}

	return dc, nil
}
