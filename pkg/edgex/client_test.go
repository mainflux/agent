package edgex

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mainflux/mainflux/logger"
)

const expectedResponse = "Response"

func TestPushOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedURL := "/operation"
		if r.URL.String() != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, r.URL.String())
		}

		expectedBody := `{"action":"start","services":["service1","service2"]}`
		bodyBytes, _ := io.ReadAll(r.Body)
		if string(bodyBytes) != expectedBody {
			t.Errorf("Expected request body %s, got %s", expectedBody, string(bodyBytes))
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(expectedResponse)); err != nil {
			t.Errorf("error writing response %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", logger.NewMock())

	response, err := client.PushOperation([]string{"start", "service1", "service2"})
	if err != nil {
		t.Errorf("Error calling PushOperation: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, response)
	}
}

func TestFetchConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedURL := "/config/start,service1,service2"
		if r.URL.String() != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, r.URL.String())
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(expectedResponse)); err != nil {
			t.Errorf("error writing response %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", logger.NewMock())

	response, err := client.FetchConfig([]string{"start", "service1", "service2"})
	if err != nil {
		t.Errorf("Error calling FetchConfig: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, response)
	}
}

func TestFetchMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedURL := "/metrics/start,service1,service2"
		if r.URL.String() != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, r.URL.String())
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(expectedResponse)); err != nil {
			t.Errorf("error writing response %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", logger.NewMock())

	response, err := client.FetchMetrics([]string{"start", "service1", "service2"})
	if err != nil {
		t.Errorf("Error calling FetchMetrics: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, response)
	}
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedURL := "/ping"
		if r.URL.String() != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, r.URL.String())
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(expectedResponse)); err != nil {
			t.Errorf("error writing response %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", logger.NewMock())

	response, err := client.Ping()
	if err != nil {
		t.Errorf("Error calling Ping: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, response)
	}
}
