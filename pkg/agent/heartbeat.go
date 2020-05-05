package agent

import (
	"sync"
	"time"
)

const (
	timeout  = 3
	interval = 10000

	online  = "online"
	offline = "offline"

	service = "service"
	device  = "device"
)

// svc keeps info on service live status.
// Services send heartbeat to nats thus updating last seen.
// When service doesnt send heartbeat for some time gets marked offline.
type svc struct {
	info    Info
	counter int
	ticker  *time.Ticker
	mu      sync.Mutex
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

func NewHeartbeat(name, svctype string) Heartbeat {
	ticker := time.NewTicker(interval * time.Millisecond)
	s := svc{
		info: Info{
			Name:     name,
			Status:   online,
			Type:     svctype,
			LastSeen: time.Now(),
		},
		counter: timeout,
		ticker:  ticker,
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
				s.counter = s.counter - 1
				if s.counter == 0 {
					s.info.Status = offline
					s.counter = timeout
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
	s.counter = timeout
	s.info.Status = online
}

func (s *svc) Info() Info {
	return s.info
}
