package services

import (
	"sync"
	"time"
)

type Status string

const (
	timeout  = 3
	interval = 10000

	Online  Status = "online"
	Offline Status = "offline"
)

type Service struct {
	Name     string
	LastSeen time.Time
	Status   Status

	counter int
	done    chan bool
	ticker  *time.Ticker
	mu      sync.Mutex
}

func NewService(name string) *Service {
	ticker := time.NewTicker(interval * time.Millisecond)
	done := make(chan bool)
	s := Service{Name: name, Status: Online, done: done, counter: timeout, ticker: ticker}
	s.Listen()
	return &s
}

func (s *Service) Listen() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				// TODO - we can disable ticker when the status gets OFFLINE
				// and on the next heartbeat enable it again
				s.mu.Lock()
				s.counter = s.counter - 1
				if s.counter == 0 {
					s.Status = Offline
					s.counter = timeout
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Service) Update() {
	s.LastSeen = time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = timeout
	s.Status = Online
}
