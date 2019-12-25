package register

import (
	"strings"

	nats "github.com/nats-io/go-nats"
)

const (
	HEARTBEAT = "heartbeat.*"
)

var _ Service = (*register)(nil)

type Service interface {
	Applications() (map[string]*Application, error)
}

type register struct {
	nc   *nats.Conn
	apps map[string]*Application
}

func New(nc *nats.Conn) (Service, error) {
	r := register{
		nc:   nc,
		apps: make(map[string]*Application),
	}

	_, err := r.nc.Subscribe(HEARTBEAT, func(msg *nats.Msg) {
		sub := msg.Subject
		appname := strings.Split(sub, ".")[1]
		// Service name is extracted from the subtopic
		// if there is multiple instances of the same service
		// we will have to add another distinction
		if _, ok := r.apps[appname]; !ok {
			a := NewApplication(appname)
			r.apps[appname] = a
			return
		}
		a := r.apps[appname]
		a.Update()
	})
	if err != nil {
		return &r, err
	}
	return &r, nil
}

func (r *register) Applications() (map[string]*Application, error) { return r.apps, nil }
