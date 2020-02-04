// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/mainflux/agent/internal/app/agent"
	"github.com/mainflux/agent/internal/pkg/config"
	export "github.com/mainflux/export/pkg/config"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/things"
)

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
	ExportConfig export.Config `json:"export_config"  mapstructure:"export_config"`
}

// Bootstrap - Retrieve device config
func Bootstrap(cfg Config, logger log.Logger, file string) error {
	retries, err := strconv.ParseUint(cfg.Retries, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid BOOTSTRAP_RETRIES value: %s", err)
	}

	retryDelaySec, err := strconv.ParseUint(cfg.RetryDelaySec, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid BOOTSTRAP_RETRY_DELAY_SECONDS value: %s", err)
	}

	logger.Info(fmt.Sprintf("Requesting config for %s from %s", cfg.ID, cfg.URL))

	dc := deviceConfig{}
	for i := 0; i < int(retries); i++ {
		dc, err = getConfig(cfg.ID, cfg.Key, cfg.URL)
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

	ic := infraConfig{}
	fmt.Println(string(dc.Content))
	if err := json.Unmarshal([]byte(dc.Content), &ic); err != nil {
		return err
	}

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

	tc := config.ThingConf{
		ID:  dc.MainfluxID,
		Key: dc.MainfluxKey,
	}
	cc := config.ChanConf{
		Control: ctrlChan,
		Data:    dataChan,
	}
	ec := config.EdgexConf{URL: ic.EdgexURL}
	lc := config.LogConf{Level: ic.LogLevel}
	mc := config.MQTTConf{URL: ic.MqttURL}

	c := config.New(sc, tc, cc, ec, lc, mc, file)

	return c.Save()
}

func getConfig(bsID, bsKey, bsSvrURL string) (deviceConfig, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", bsSvrURL, bsID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return deviceConfig{}, err
	}

	req.Header.Add("Authorization", bsKey)
	resp, err := client.Do(req)
	if err != nil {
		return deviceConfig{}, err
	}
	if resp.StatusCode == http.StatusForbidden {
		return deviceConfig{}, errors.New("Unauthorized access")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return deviceConfig{}, err
	}
	defer resp.Body.Close()

	dc := deviceConfig{}
	if err := json.Unmarshal(body, &dc); err != nil {
		return deviceConfig{}, err
	}

	return dc, nil
}
