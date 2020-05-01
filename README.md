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

```bash
go get github.com/mainflux/agent
cd $GOPATH/github.com/mainflux/agent
```

Make:
```bash
make
```

## Usage
Get Nats server and start it, by default it starts on port `4222`

```bash
go get github.com/nats-io/gnatsd
gnatsd
```

Create gateway configuration with [Provision][provision] service or through [Mainflux UI][mfxui].  

Start Agent with:

```bash
MF_AGENT_BOOTSTRAP_ID=<bootstrap_id> \
MF_AGENT_BOOTSTRAP_KEY=<bootstrap_key> \
MF_AGENT_BOOTSTRAP_URL=https://mainflux.com/bs/things/bootstrap \
build/mainflux-agent
```

### Config
Agent configuration is kept in `config.toml` if not otherwise specified with env var.

Example configuration:
```toml
[Agent]

  [Agent.channels]
    control = ""
    data = ""

  [Agent.edgex]
    url = "http://localhost:48090/api/v1/"

  [Agent.log]
    level = "info"

  [Agent.mqtt]
    ca_path = "ca.crt"
    cert_path = "thing.crt"
    mtls = false
    password = ""
    priv_key_path = "thin.key"
    qos = 0
    retain = false
    skip_tls_ver = false
    url = "localhost:1883"
    username = ""

  [Agent.server]
    nats_url = "localhost:4222"
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
(i.e. app needs to PUB/SUB on `/channels/<control_channel_id>/messages/req` and `/channels/<control_channel_id>/messages/res`).

## Sending commands to other services
You can send commands to other services that are subscribed on the same Nats server as Agent.  
Commands are being sent via MQTT to topic:   
* `channels/<control_channel_id>/messages/services/<service_name>/<subtopic>`  

when messages is received Agent forwards them to Nats on subject:   
* `commands.<service_name>.<subtopic>`.  

Payload is up to the application and service itself.

Example of on command can be:

```bash
mosquitto_pub -u <thing_id> -P <thing_key> -t channels/<control_channel_id>/messages/services/adc -h <mqtt_host> -p 1883  -m  "[{\"bn\":\"1:\", \"n\":\"read\", \"vs\":\"temperature\"}]"
```

## Heartbeat service
Services running on the same host can publish to `heartbeat.<service-name>.<service-type>` a heartbeat message.  
Agent will keep a record on those service and update their `live` status. If heartbeat is not received in 10 sec it marks it `offline`.
Upon next heartbeat service will be marked `online` again.
To check services that are currently registered to agent you can:

```bash
curl -s -S X GET http://localhost:9000/services
```
```json
[
  {
    "Name": "duster",
    "LastSeen": "2020-04-28T18:06:56.158130519+02:00",
    "Status": "offline",
    "Type": "test",
    "Terminal": 0
  },
  {
    "Name": "scrape",
    "LastSeen": "2020-04-28T18:06:39.58849766+02:00",
    "Status": "offline",
    "Type": "test",
    "Terminal": 0
  }
]
```

Or you can send a command via MQTT to Agent and receive response on MQTT topic like this:

In one terminal subscribe for result:

```bash
mosquitto_sub -u <thing_id> -P <thing_key> -t channels/<control_channel_id>/messages/req -h <mqtt_host> -p 1883  
```

In another terminal publish request to view the list of services:

```bash
mosquitto_pub -u <thing_id> -P <thing_key> -t channels/<control_channel_id>/messages/req -h <mqtt_host> -p 1883  -m  '[{"bn":"1:", "n":"config", "vs":"view"}]'  
```

Check the output in terminal where you subscribed for results. You should see something like:

```json
[
  {
    "bn": "1",
    "n": "view",
    "t": 1588091188.8872917,
    "vs": "[{\"Name\":\"duster\",\"LastSeen\":\"2020-04-28T18:06:56.158130519+02:00\",\"Status\":\"offline\",\"Type\":\"test\",\"Terminal\":0},{\"Name\":\"scrape\",\"LastSeen\":\"2020-04-28T18:06:39.58849766+02:00\",\"Status\":\"offline\",\"Type\":\"test\",\"Terminal\":0}]"                                                                       
  }
]
```


## How to save config via agent
Agent can be used to send configuration file for the [Export][export] service from cloud to gateway via MQTT.
Here is the example command

```bash
mosquitto_pub -u <thing_id> -P <thing_key> -t channels/<control_ch_id>/messages/req -h localhost -p 18831  -m  "[{\"bn\":\"1:\", \"n\":\"config\", \"vs\":\"config.toml, RmlsZSA9ICIuLi9jb25maWdzL2NvbmZpZy50b21sIgoKW2V4cF0KICBsb2dfbGV2ZWwgPSAiZGVidWciCiAgbmF0cyA9ICJuYXRzOi8vMTI3LjAuMC4xOjQyMjIiCiAgcG9ydCA9ICI4MTcwIgoKW21xdHRdCiAgY2FfcGF0aCA9ICJjYS5jcnQiCiAgY2VydF9wYXRoID0gInRoaW5nLmNydCIKICBjaGFubmVsID0gIiIKICBob3N0ID0gInRjcDovL2xvY2FsaG9zdDoxODgzIgogIG10bHMgPSBmYWxzZQogIHBhc3N3b3JkID0gImFjNmI1N2UwLTliNzAtNDVkNi05NGM4LWU2N2FjOTA4NjE2NSIKICBwcml2X2tleV9wYXRoID0gInRoaW5nLmtleSIKICBxb3MgPSAwCiAgcmV0YWluID0gZmFsc2UKICBza2lwX3Rsc192ZXIgPSBmYWxzZQogIHVzZXJuYW1lID0gIjRhNDM3ZjQ2LWRhN2ItNDQ2OS05NmI3LWJlNzU0YjVlOGQzNiIKCltbcm91dGVzXV0KICBtcXR0X3RvcGljID0gIjRjNjZhNzg1LTE5MDAtNDg0NC04Y2FhLTU2ZmI4Y2ZkNjFlYiIKICBuYXRzX3RvcGljID0gIioiCg==\"}]"
```

`vs="config_file_path, file_cont_base64"` - vs determines where to save file and contains file content in base64 encoding payload.  

Here is an example how to make payload for the command:
```go
b,_ := toml.Marshal(export.Config)
payload := base64.StdEncoding.EncodeToString(b)
```

## License

[Apache-2.0](LICENSE)

[grc-badge]: https://goreportcard.com/badge/github.com/mainflux/agent
[grc-url]: https://goreportcard.com/report/github.com/mainflux/agent
[docs]: http://mainflux.readthedocs.io
[gitter]: https://gitter.im/mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge
[gitter-badge]: https://badges.gitter.im/Join%20Chat.svg
[license]: https://img.shields.io/badge/license-Apache%20v2.0-blue.svg

[export]: https://github.com/mainflux/export
[provision]: https://github.com/mainflux/mainflux/tree/master/provision
[mfxui]: https://github.com/mainflux/ui
