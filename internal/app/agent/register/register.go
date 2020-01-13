package register

import (
	"fmt"
	"strings"

	"github.com/mainflux/mainflux/logger"
	nats "github.com/nats-io/go-nats"
)

const (
	Hearbeat = "heartbeat.*"
)

var _ Service = (*register)(nil)

type Service interface {
	Applications() map[string]*Application
}

type register struct {
	nc     *nats.Conn
	apps   map[string]*Application
	logger logger.Logger
}

func New(nc *nats.Conn, l logger.Logger) (Service, error) {
	r := register{
		nc:     nc,
		apps:   make(map[string]*Application),
		logger: l,
	}

	_, err := r.nc.Subscribe(Hearbeat, func(msg *nats.Msg) {

		sub := msg.Subject
		tok := strings.Split(sub, ".")
		if len(tok) < 2 {
			l.Error(fmt.Sprintf("Failed: Subject has incorrect length %s" + sub))
			return
		}
		appname := tok[1]
		// Service name is extracted from the subtopic
		// if there is multiple instances of the same service
		// we will have to add another distinction
		if _, ok := r.apps[appname]; !ok {
			a := NewApplication(appname)
			r.apps[appname] = a
			l.Info(fmt.Sprintf("Application '%s' registered", appname))
		}
		a := r.apps[appname]
		a.Update()
	})
	if err != nil {
		return &r, err
	}
	return &r, nil
}

func (r *register) Applications() map[string]*Application { return r.apps }
