package register

import (
	"sync"
	"time"
)

type Status string

const (
	timeout  = 10
	interval = 10000

	ONLINE  Status = "online"
	OFFLINE Status = "offline"
)

type Application struct {
	Name string
	// RegName  string
	// Num      int
	LastSeen time.Time
	Status   Status

	counter int
	done    *chan bool
	ticker  *time.Ticker
	mu      sync.Mutex
}

func NewApplication(name string) *Application {
	ticker := time.NewTicker(interval * time.Millisecond)
	done := new(chan bool)
	a := Application{Name: name, done: done, ticker: ticker}
	go func() {
		for {
			select {
			case <-*a.done:
				return
			case <-ticker.C:
				a.mu.Lock()
				a.counter = a.counter - 1
				if a.counter == 0 {
					a.Status = OFFLINE
				}
				a.mu.Unlock()
			}
		}
	}()
	return &a
}

func (a *Application) Update() {
	a.LastSeen = time.Now()
	a.mu.Lock()
	defer a.mu.Unlock()
	a.counter = timeout
	a.Status = ONLINE
}

func (a *Application) Done() {
	*a.done <- true
}
