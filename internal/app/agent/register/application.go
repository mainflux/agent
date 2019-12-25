package register

import (
	"time"
)

type Status string

const (
	ONLINE  Status = "online"
	OFFLINE Status = "offline"
)

type Application struct {
	Name     string
	RegName  string
	Num      int
	LastSeen time.Time
	Status   Status
}
