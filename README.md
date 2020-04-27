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
```
go get github.com/nats-io/gnatsd
gnatsd
MF_AGENT_BOOTSTRAP_ID=11:22:33:44 MF_AGENT_BOOTSTRAP_KEY=5341a625-53a5-4a25-be74-5b648243fcbc MF_AGENT_NATS_URL=nats://localhost:4223 MF_AGENT_BOOTSTRAP_URL=https://mainflux.com/bs/things/bootstrap build/mainflux-agent
```

### Config
Agent configuration is kept in `config.toml` if not otherwise specified with env var.

Example configuration:
```
File = "config.toml"

[Agent]

  [Agent.channels]
    control = "1280c838-df6f-4532-b3f0-608d7c6e00ef"
    data = "5e5f1397-4c42-40f1-84f0-72cc7c7efa5a"

  [Agent.edgex]
    url = "http://localhost:48090/api/v1/"

  [Agent.log]
    level = "info"

  [Agent.mqtt]
    ca_path = ""
    cert_path = ""
    mtls = false
    password = "0c6fd027-79fc-4fca-be27-814f7e842cf9"
    priv_key_path = ""
    qos = 0
    retain = false
    skip_tls_ver = false
    url = "84.201.171.65:1883"
    username = "fc430ab3-6765-49b0-af70-8629fe5d3900"

  [Agent.server]
    nats_url = "localhost:4223"
    port = "9000"

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

