package terminal

import (
	"errors"
	"testing"
	"time"

	"github.com/mainflux/mainflux/logger"
	"github.com/stretchr/testify/assert"
)

const (
	uuid    = "test-uuid"
	timeout = 5 * time.Second
)

// MockPublish is a mock function for the publish function used in NewSession.
func mockPublish(channel, payload string) error {
	return nil
}

func mockPublishFail(channel, payload string) error {
	return errors.New("")
}

func TestWrite(t *testing.T) {
	t.Run("successful publish", func(t *testing.T) {
		session, err := NewSession(uuid, timeout, mockPublish, logger.NewMock())
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		// Simulate writing data to the session
		data := []byte("test data")
		n, err := session.Write(data)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		assert.Equal(t, len(data), n)
	})
	t.Run("failed publish", func(t *testing.T) {
		session, err := NewSession(uuid, timeout, mockPublishFail, logger.NewMock())
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		// Simulate writing data to the session
		data := []byte("test data")
		_, err = session.Write(data)
		assert.NotNil(t, err)
	})
}

func TestSend(t *testing.T) {
	session, err := NewSession(uuid, timeout, mockPublish, logger.NewMock())
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Simulate sending data to the session
	data := []byte("test data")

	if err = session.Send(data); err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

}

func TestIsDone(t *testing.T) {
	publish := mockPublish

	session, err := NewSession(uuid, timeout, publish, logger.NewMock())
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Wait for the "done" channel to be closed or for a timeout, and perform assertions accordingly.
	select {
	case <-session.IsDone():
		// Session is done as expected.
	case <-time.After(10 * time.Second):
		t.Fatalf("Expected session to be done, but it is still running.")
	}
}
