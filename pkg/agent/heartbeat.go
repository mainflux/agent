package agent

import (
	"sync"
	"time"
)

const (
	online  = "online"
	offline = "offline"

	service = "service"
	device  = "device"
)

// svc keeps info on service live status.
// Services send heartbeat to nats thus updating last seen.
// When service doesnt send heartbeat for some time gets marked offline.
type svc struct {
	info        Info
	counter     int
	numInterval int
	ticker      *time.Ticker
	mu          sync.Mutex
}

type Info struct {
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
	Status   string    `json:"status"`
	Type     string    `json:"type"`
	Terminal int       `json:"terminal"`
}

// Heartbeat specifies api for updating status and keeping track on services
// that are sending heartbeat to NATS.
type Heartbeat interface {
	Update()
	Info() Info
}

// interval - duration of interval
// if service doesnt send heartbeat during  interval it is marked offline
func NewHeartbeat(name, svcType string, interval time.Duration) Heartbeat {
	ticker := time.NewTicker(interval)
	s := svc{
		info: Info{
			Name:     name,
			Status:   online,
			Type:     svcType,
			LastSeen: time.Now(),
		},
		ticker: ticker,
	}
	s.listen()
	return &s
}

func (s *svc) listen() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				// TODO - we can disable ticker when the status gets OFFLINE
				// and on the next heartbeat enable it again
				s.mu.Lock()
				if s.counter == 0 {
					s.info.Status = offline
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *svc) Update() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.info.LastSeen = time.Now()
	s.counter = 1
	s.info.Status = online
}

func (s *svc) Info() Info {
	return s.info
}
