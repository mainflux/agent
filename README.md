# Mainflux IoT Agent

[![build][ci-badge]][ci-url]
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

Here `thing` is a Mainflux thing, and control channel from `channels` is used with `req` and `res` subtopic
(i.e. app needs to PUB/SUB on `/channels/<channel_id>/messages/req` and `/channels/<channel_id>/messages/res`).

## License

[Apache-2.0](LICENSE)

[ci-badge]: https://semaphoreci.com/api/v1/mainflux/agent/branches/master/badge.svg
[ci-url]: https://semaphoreci.com/mainflux/agent
[docs]: http://mainflux.readthedocs.io
[gitter]: https://gitter.im/mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge
[gitter-badge]: https://badges.gitter.im/Join%20Chat.svg
[license]: https://img.shields.io/badge/license-Apache%20v2.0-blue.svg

