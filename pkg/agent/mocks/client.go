// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mocks

// mockClient - holds data for Edgex mockClient
type mockClient struct {
}

// NewmockClient - Creates ne EdgeX mockClient
func NewEdgexClient() *mockClient {
	return &mockClient{}
}

// PushOperation - pushes operation to EdgeX components
func (ec *mockClient) PushOperation(cmdArr []string) (string, error) {
	return string("body"), nil
}

// FetchConfig - fetches config from EdgeX components
func (ec *mockClient) FetchConfig(cmdArr []string) (string, error) {
	return string("body"), nil
}

// FetchMetrics - fetches metrics from EdgeX components
func (ec *mockClient) FetchMetrics(cmdArr []string) (string, error) {
	return string("body"), nil
}

// Ping - ping EdgeX SMA
func (ec *mockClient) Ping() (string, error) {
	return string("body"), nil
}
