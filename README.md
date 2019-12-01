# Mainflux IoT Agent
![agent](./docs/img/agent.png =250x)

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

`thing` is a Mainflux thing, and control channel from `channels` is used with `req` and `res` subtopic
(i.e. app needgs to pub/sub on `/channels/<channel_id>/messages/req` and `/channels/<channel_id>/messages/res`).



