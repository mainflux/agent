package agent

import (
	"testing"
	"time"
)

func TestNewHeartbeat(t *testing.T) {
	name := "TestService"
	serviceType := "TestType"
	interval := 2 * time.Second

	heartbeat := NewHeartbeat(name, serviceType, interval)

	// Check initial status and info
	info := heartbeat.Info()
	if info.Name != name {
		t.Errorf("Expected name to be %s, but got %s", name, info.Name)
	}
	if info.Type != serviceType {
		t.Errorf("Expected type to be %s, but got %s", serviceType, info.Type)
	}
	if info.Status != online {
		t.Errorf("Expected initial status to be %s, but got %s", online, info.Status)
	}
}

func TestHeartbeat_Update(t *testing.T) {
	name := "TestService"
	serviceType := "TestType"
	interval := 2 * time.Second

	heartbeat := NewHeartbeat(name, serviceType, interval)

	// Sleep for more than the interval to simulate an update
	time.Sleep(3 * time.Second)

	heartbeat.Update()

	// Check if the status has been updated to online
	info := heartbeat.Info()
	if info.Status != online {
		t.Errorf("Expected status to be %s, but got %s", online, info.Status)
	}
}

func TestHeartbeat_StatusOffline(t *testing.T) {
	name := "TestService"
	serviceType := "TestType"
	interval := 2 * time.Second

	heartbeat := NewHeartbeat(name, serviceType, interval)

	// Sleep for more than two intervals to simulate offline status
	time.Sleep(5 * time.Second)

	// Check if the status has been updated to offline
	info := heartbeat.Info()
	if info.Status != offline {
		t.Errorf("Expected status to be %s, but got %s", offline, info.Status)
	}
}
