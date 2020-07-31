module github.com/mainflux/agent

go 1.14

require (
	github.com/creack/pty v1.1.9
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/edgexfoundry/go-mod-core-contracts v0.1.48
	github.com/go-kit/kit v0.9.0
	github.com/go-zoo/bone v1.3.0
	github.com/mainflux/export v0.0.0-20200323141637-120ec7179230
	github.com/mainflux/mainflux v0.0.0-20200212173448-51cd0524a11e
	github.com/mainflux/senml v1.0.0
	github.com/nats-io/nats.go v1.9.1
	github.com/pelletier/go-toml v1.7.0
	github.com/prometheus/client_golang v1.4.1
	github.com/stretchr/testify v1.4.0
	robpike.io/filter v0.0.0-20150108201509-2984852a2183
)

replace github.com/mainflux/export => ../export
