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
	a := Application{Name: name, Status: ONLINE, done: done, ticker: ticker}
	a.Listen()
	return &a
}

func (a *Application) Listen() {
	go func() {
		for {
			select {
			case <-*a.done:
				return
			case <-a.ticker.C:
				// TODO - we can disable ticker when the status gets OFFLINE
				// and on the next heartbeat enable it again
				a.mu.Lock()
				a.counter = a.counter - 1
				if a.counter == 0 {
					a.Status = OFFLINE
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
	a.Status = ONLINE
}

func (a *Application) Done() {
	*a.done <- true
}
