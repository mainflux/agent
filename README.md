# Mainflux IoT Agent

![](https://github.com/mainflux/agent/workflows/Go/badge.svg)
[![go report card][grc-badge]][grc-url]
[![license][license]](LICENSE)
[![chat][gitter-badge]][gitter]

<p align="center">
  <img width="30%" height="30%" src="./docs/img/agent.png">
</p>

Mainflux IoT Agent is a communication, execution and SW management agent for Mainflux system.

## Install
Get the code:

```
go get github.com/mainflux/agent
cd $GOPATH/github.com/mainflux/agent
```

Make:
```
make
```

## Usage
### Config
Agent configuration is kept in `cmd/config.toml`.

Example configuration:
```
[server]
port = "9000"

[thing]
id = "90f9aff0-73f6-476b-a4da-d277ab99d3ce"
key = "e606945e-47bd-405f-85f3-a4c55bd069c8"

[channels]
control = "3ace3fa3-aa84-4a02-b0ab-6d594268dc77"
data = "0bdeea73-8de1-420c-8c4d-a4e7d6c46e3b"

[edgex]
url = "http://localhost:48090/api/v1/"

[log]
level = "info"

[mqtt]
url = "localhost:1883"
```

Environment:
| Variable                               | Description                                                   | Default                           |
|----------------------------------------|---------------------------------------------------------------|-----------------------------------|
|	MF_AGENT_CONFIG_FILE                   | Location of configuration file                                | config.toml                       |
|	MF_AGENT_LOG_LEVEL                     | Log level                                                     | info                              |
|	MF_AGENT_EDGEX_URL                     | Edgex base url                                                | http://localhost:48090/api/v1/    |
|	MF_AGENT_MQTT_URL                      | MQTT broker url                                               | localhost:1883                    |
|	MF_AGENT_HTTP_PORT                     | Agent http port                                               | 9000                              |
|	MF_AGENT_BOOTSTRAP_URL                 | Mainflux bootstrap url                                        | http://localhost:8202/things/bootstrap|
|	MF_AGENT_BOOTSTRAP_ID                  | Mainflux bootstrap id                                         |                                   |
|	MF_AGENT_BOOTSTRAP_KEY                 | Mainflux boostrap key                                         |                                   |
|	MF_AGENT_BOOTSTRAP_RETRIES             | Number of retries for bootstrap procedure                     | 5                                 |
|	MF_AGENT_BOOTSTRAP_RETRY_DELAY_SECONDS | Number of seconds between retries                             | 10                                |
|	MF_AGENT_CONTROL_CHANNEL               | Channel for sending controls, commands                        |                                   |
|	MF_AGENT_DATA_CHANNEL                  | Channel for data sending                                      |                                   |
|	MF_AGENT_ENCRYPTION                    | Encryption                                                    | false                             |
|	MF_AGENT_NATS_URL                      | Nats url                                                      | nats://localhost:4222             |
|	MF_AGENT_MQTT_USERNAME                 | MQTT username, Mainflux thing id                              |                                   |
|	MF_AGENT_MQTT_PASSWORD                 | MQTT password, Mainflux thing key                             |                                   |
|	MF_AGENT_MQTT_SKIP_TLS                 | Skip TLS verification                                         | true                              |
|	MF_AGENT_MQTT_MTLS                     | Use MTLS for MQTT                                             | false                             |
|	MF_AGENT_MQTT_CA                       | Location for CA certificate for MTLS                          | ca.crt                            |
|	MF_AGENT_MQTT_QOS                      | QoS                                                           | 0                                 |
|	MF_AGENT_MQTT_RETAIN                   | MQTT retain                                                   | false                             |
|	MF_AGENT_MQTT_CLIENT_CERT              | Location of client certificate for MTLS                       | thing.cert                        |
|	MF_AGENT_MQTT_CLIENT_PK                | Location of client certificate key for MTLS                   | thing.key                         |

Here `thing` is a Mainflux thing, and control channel from `channels` is used with `req` and `res` subtopic
(i.e. app needs to PUB/SUB on `/channels/<channel_id>/messages/req` and `/channels/<channel_id>/messages/res`).

## License

[Apache-2.0](LICENSE)

[grc-badge]: https://goreportcard.com/badge/github.com/mainflux/agent
[grc-url]: https://goreportcard.com/report/github.com/mainflux/agent
[docs]: http://mainflux.readthedocs.io
[gitter]: https://gitter.im/mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge
[gitter-badge]: https://badges.gitter.im/Join%20Chat.svg
[license]: https://img.shields.io/badge/license-Apache%20v2.0-blue.svg

