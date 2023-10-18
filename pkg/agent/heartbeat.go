package agent

import (
	"context"
	"sync"
	"time"
)

const (
	online  = "online"
	offline = "offline"
)

// svc keeps info on service live status.
// Services send heartbeat to nats thus updating last seen.
// When service doesnt send heartbeat for some time gets marked offline.
type svc struct {
	info     Info
	interval time.Duration
	ticker   *time.Ticker
	mu       sync.Mutex
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
	Close()
}

// interval - duration of interval
// if service doesnt send heartbeat during  interval it is marked offline.
func NewHeartbeat(ctx context.Context, name, svcType string, interval time.Duration) Heartbeat {
	ticker := time.NewTicker(interval)
	s := svc{
		info: Info{
			Name:     name,
			Status:   online,
			Type:     svcType,
			LastSeen: time.Now(),
		},
		ticker:   ticker,
		interval: interval,
	}
	go s.listen(ctx)
	return &s
}

func (s *svc) listen(ctx context.Context) {
	for {
		select {
		case <-s.ticker.C:
			// TODO - we can disable ticker when the status gets OFFLINE
			// and on the next heartbeat enable it again.
			s.mu.Lock()
			if time.Now().After(s.info.LastSeen.Add(s.interval)) {
				s.info.Status = offline
			}
			s.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (s *svc) Update() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.info.LastSeen = time.Now()
	s.info.Status = online
}

func (s *svc) Info() Info {
	s.mu.Lock()
	defer s.mu.Unlock()
	info := s.info
	return info
}

func (s *svc) Close() {
	s.ticker.Stop()
}
