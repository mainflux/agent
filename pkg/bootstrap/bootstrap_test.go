package bootstrap

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mainflux/agent/pkg/agent"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestBootstrap(t *testing.T) {
	// Create a mock HTTP server to handle requests from the getConfig function.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Thing mockKey" && r.Header.Get("Authorization") != "Thing invalidChannels" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Authorization") == "Thing invalidChannels" {
			// Simulate a malformed response.
			resp := `
			{
				"thing_id": "e22c383a-d2ab-47c1-89cd-903955da993d",
				"thing_key": "fc987711-1828-461b-aa4b-16d5b2c642fe",
				"channels": [
				  {
					"id": "fa5f9ba8-a1fc-4380-9edb-d0c23eaa24ec",
					"name": "control-channel",
					"metadata": {
					  "type": "control"
					}
				  }
				],
				"content": "{\"agent\":{\"edgex\":{\"url\":\"http://localhost:48090/api/v1/\"},\"heartbeat\":{\"interval\":\"30s\"},\"log\":{\"level\":\"debug\"},\"mqtt\":{\"mtls\":false,\"qos\":0,\"retain\":false,\"skip_tls_ver\":true,\"url\":\"tcp://mainflux-domain.com:1883\"},\"server\":{\"nats_url\":\"localhost:4222\",\"port\":\"9000\"},\"terminal\":{\"session_timeout\":\"30s\"}},\"export\":{\"exp\":{\"cache_db\":\"0\",\"cache_pass\":\"\",\"cache_url\":\"localhost:6379\",\"log_level\":\"debug\",\"nats\":\"nats://localhost:4222\",\"port\":\"8172\"},\"mqtt\":{\"ca_path\":\"ca.crt\",\"cert_path\":\"thing.crt\",\"channel\":\"\",\"host\":\"tcp://mainflux-domain.com:1883\",\"mtls\":false,\"password\":\"\",\"priv_key_path\":\"thing.key\",\"qos\":0,\"retain\":false,\"skip_tls_ver\":false,\"username\":\"\"},\"routes\":[{\"mqtt_topic\":\"\",\"nats_topic\":\"channels\",\"subtopic\":\"\",\"type\":\"mfx\",\"workers\":10},{\"mqtt_topic\":\"\",\"nats_topic\":\"export\",\"subtopic\":\"\",\"type\":\"default\",\"workers\":10}]}}"
			  }
			`
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, resp)
			return
		}
		// Simulate a successful response.
		resp := `
		{
			"thing_id": "e22c383a-d2ab-47c1-89cd-903955da993d",
			"thing_key": "fc987711-1828-461b-aa4b-16d5b2c642fe",
			"channels": [
			  {
				"id": "fa5f9ba8-a1fc-4380-9edb-d0c23eaa24ec",
				"name": "control-channel",
				"metadata": {
				  "type": "control"
				}
			  },
			  {
				"id": "24e5473e-3cbe-43d9-8a8b-a725ff918c0e",
				"name": "data-channel",
				"metadata": {
				  "type": "data"
				}
			  },
			  {
				"id": "1eac45c2-0f72-4089-b255-ebd2e5732bbb",
				"name": "export-channel",
				"metadata": {
				  "type": "export"
				}
			  }
			],
			"content": "{\"agent\":{\"edgex\":{\"url\":\"http://localhost:48090/api/v1/\"},\"heartbeat\":{\"interval\":\"30s\"},\"log\":{\"level\":\"debug\"},\"mqtt\":{\"mtls\":false,\"qos\":0,\"retain\":false,\"skip_tls_ver\":true,\"url\":\"tcp://mainflux-domain.com:1883\"},\"server\":{\"nats_url\":\"localhost:4222\",\"port\":\"9000\"},\"terminal\":{\"session_timeout\":\"30s\"}},\"export\":{\"exp\":{\"cache_db\":\"0\",\"cache_pass\":\"\",\"cache_url\":\"localhost:6379\",\"log_level\":\"debug\",\"nats\":\"nats://localhost:4222\",\"port\":\"8172\"},\"mqtt\":{\"ca_path\":\"ca.crt\",\"cert_path\":\"thing.crt\",\"channel\":\"\",\"host\":\"tcp://mainflux-domain.com:1883\",\"mtls\":false,\"password\":\"\",\"priv_key_path\":\"thing.key\",\"qos\":0,\"retain\":false,\"skip_tls_ver\":false,\"username\":\"\"},\"routes\":[{\"mqtt_topic\":\"\",\"nats_topic\":\"channels\",\"subtopic\":\"\",\"type\":\"mfx\",\"workers\":10},{\"mqtt_topic\":\"\",\"nats_topic\":\"export\",\"subtopic\":\"\",\"type\":\"default\",\"workers\":10}]}}"
		  }
		`
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, resp)
	}))
	defer mockServer.Close()
	mockLogger := logger.NewMock()
	tests := []struct {
		name        string
		config      Config
		file        string
		expectedErr error
	}{
		{
			name:        "invalid retries type",
			config:      Config{Retries: "invalid"},
			expectedErr: errInvalidBootstrapRetriesValue,
		},
		{
			name:        "zero retires",
			config:      Config{Retries: "0"},
			expectedErr: nil,
		},
		{
			name:        "invalid retry delay",
			config:      Config{Retries: "1", RetryDelaySec: "e"},
			expectedErr: errInvalidBootstrapRetryDelay,
		},
		{
			name:        "authorization error",
			config:      Config{Retries: "1", RetryDelaySec: "1", URL: mockServer.URL, Key: "wrongKey"},
			expectedErr: nil,
		},
		{
			name:        "malformed channels",
			config:      Config{Retries: "1", RetryDelaySec: "1", URL: mockServer.URL, Key: "invalidChannels"},
			expectedErr: errors.ErrMalformedEntity,
		},
		{
			name:        "successful configuration",
			config:      Config{Retries: "1", RetryDelaySec: "1", URL: mockServer.URL, Key: "mockKey"},
			expectedErr: agent.ErrWritingToml,
		},
		{
			name:        "successful configuration",
			config:      Config{Retries: "1", RetryDelaySec: "1", URL: mockServer.URL, Key: "mockKey"},
			expectedErr: nil,
			file:        "config.toml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Bootstrap(test.config, mockLogger, test.file)
			assert.True(t, errors.Contains(err, test.expectedErr), fmt.Sprintf("expected %v got %v", test.expectedErr, err))
		})
	}
	// cleanup.
	t.Cleanup(func() {
		os.Remove("config.toml")
	})
}
