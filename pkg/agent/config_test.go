package agent

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary config file for testing.
	tempFile, err := os.CreateTemp("", "config.toml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile2, err := os.CreateTemp("", "invalid.toml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile2.Name())

	sampleConfig := `
	File = "config.toml"

    [channels]
      control = ""
      data = ""

    [edgex]
      url = "http://localhost:48090/api/v1/"

    [heartbeat]
      interval = "10s"

    [log]
      level = "info"

    [mqtt]
      ca_cert = ""
      ca_path = "ca.crt"
      cert_path = "thing.cert"
      client_cert = ""
      client_key = ""
      mtls = false
      password = ""
      priv_key_path = "thing.key"
      qos = 0
      retain = false
      skip_tls_ver = true
      url = "localhost:1883"
      username = ""

    [server]
      nats_url = "nats://127.0.0.1:4222"
      port = "9999"

    [terminal]
      session_timeout = "1m0s"
`

	if _, writeErr := tempFile.WriteString(sampleConfig); writeErr != nil {
		t.Fatalf("Failed to write to temporary file: %v", writeErr)
	}
	tempFile.Close()

	if _, writeErr := tempFile2.WriteString(strings.ReplaceAll(sampleConfig, "[", "")); writeErr != nil {
		t.Fatalf("Failed to write to temporary file: %v", writeErr)
	}
	tempFile2.Close()

	tests := []struct {
		name        string
		fileName    string
		expectedErr error
	}{
		{
			name:        "failed to read file",
			fileName:    "invalidFile.toml",
			expectedErr: errReadingFile,
		},
		{
			name:        "invalid toml",
			fileName:    tempFile2.Name(),
			expectedErr: errUnmarshalToml,
		},
		{
			name:        "successful read",
			fileName:    tempFile.Name(),
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ReadConfig(test.fileName)
			assert.True(t, errors.Contains(err, test.expectedErr), fmt.Sprintf("expected %v got %v", test.expectedErr, err))
		})
	}
}
