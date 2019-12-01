// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package edgex

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/mainflux/mainflux/logger"

	model "github.com/edgexfoundry/go-mod-core-contracts/models"
)

// Client - holds data for Edgex Client
type Client struct {
	url    string
	logger log.Logger
}

// NewClient - Creates ne EdgeX client
func NewClient(edgexURL string, logger log.Logger) *Client {
	return &Client{
		url:    edgexURL,
		logger: logger,
	}
}

// PushOperation - pushes operation to EdgeX components
func (ec *Client) PushOperation(cmdArr []string) (string, error) {
	url := ec.url + "operation"

	m := model.Operation{
		Action:   cmdArr[0],
		Services: cmdArr[1:],
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FetchConfig - fetches config from EdgeX components
func (ec *Client) FetchConfig(cmdArr []string) (string, error) {
	cmdStr := strings.Replace(strings.Join(cmdArr, ","), " ", "", -1)
	url := ec.url + "config/" + cmdStr

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FetchMetrics - fetches metrics from EdgeX components
func (ec *Client) FetchMetrics(cmdArr []string) (string, error) {
	cmdStr := strings.Replace(strings.Join(cmdArr, ","), " ", "", -1)
	url := ec.url + "metrics/" + cmdStr

	resp, err := http.Get(url)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Ping - ping EdgeX SMA
func (ec *Client) Ping() (string, error) {
	url := ec.url + "ping"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
