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

type svc struct {
	name     string
	lastSeen time.Time
	status   string
	typ      string

	counter int
	ticker  *time.Ticker
	mu      sync.Mutex
}

type Info struct {
	Name     string
	LastSeen time.Time
	Status   string
	Type     string
	Terminal int
}

// Heartbeat specifies api for updating status and keeping track on services
// that are sending heartbeat to NATS.
type Heartbeat interface {
	Update()
	Info() Info
}

func NewHeartbeat(name, svctype string) Heartbeat {
	ticker := time.NewTicker(interval * time.Millisecond)
	s := svc{name: name, status: online, typ: svctype, counter: timeout, ticker: ticker}
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
					s.status = offline
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
	s.lastSeen = time.Now()
	s.counter = timeout
	s.status = online
}

func (s *svc) Info() Info {
	return Info{
		Name:     s.name,
		LastSeen: s.lastSeen,
		Status:   s.status,
		Type:     s.typ,
	}
}
