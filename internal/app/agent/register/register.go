package register

import (
	"fmt"
	"strings"
	"time"

	nats "github.com/nats-io/go-nats"
)

const (
	HEARTBEAT = "heartbeat.*"
)

var _ Service = (*register)(nil)

type Service interface {
	Add(a Application) error
	Remove(a Application) error
	List() ([]Application, error)
}

type register struct {
	nc   *nats.Conn
	apps map[string]map[string]Application
}

func New(nc *nats.Conn) (Service, error) {
	r := register{
		nc:   nc,
		apps: make(map[string]map[string]Application),
	}

	_, err := r.nc.Subscribe(HEARTBEAT, func(msg *nats.Msg) {
		sub := msg.Subject
		appname := strings.Split(sub, ".")[1]
		if amap, ok := r.apps[appname]; !ok {
			r.apps[appname] = make(map[string]Application)
			num := 0
			regname := fmt.Sprintf("%s-%d", appname, num)
			a := Application{
				Name:     appname,
				RegName:  regname,
				Num:      num,
				LastSeen: time.Now(),
				Status:   ONLINE,
			}
			r.apps[appname][regname] = a
			return
		}
		amap[appname]

	})
	if err != nil {
		return r, err
	}

	return r, nil
}

func (r register) Add(a Application) error      { return nil }
func (r register) Remove(a Application) error   { return nil }
func (r register) List() ([]Application, error) { return []Application{}, nil }
