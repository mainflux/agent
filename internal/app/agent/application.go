package agent

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

type Application struct {
	Name     string
	LastSeen time.Time
	Status   Status

	counter int
	done    chan bool
	ticker  *time.Ticker
	mu      sync.Mutex
}

func NewApplication(name string) *Application {
	ticker := time.NewTicker(interval * time.Millisecond)
	done := make(chan bool)
	a := Application{Name: name, Status: Online, done: done, counter: timeout, ticker: ticker}
	a.Listen()
	return &a
}

func (a *Application) Listen() {
	go func() {
		for {
			select {
			case <-a.ticker.C:
				// TODO - we can disable ticker when the status gets OFFLINE
				// and on the next heartbeat enable it again
				a.mu.Lock()
				a.counter = a.counter - 1
				if a.counter == 0 {
					a.Status = Offline
					a.counter = timeout
				}
				a.mu.Unlock()
			}
		}
	}()
}

func (a *Application) Update() {
	a.LastSeen = time.Now()
	a.mu.Lock()
	defer a.mu.Unlock()
	a.counter = timeout
	a.Status = Online
}
