# Message Transformers

A transformer service consumes events published by Mainflux adapters (such as MQTT and HTTP adapters) and transforms them to an arbitrary message format. A transformer can be imported as a standalone package and used for message transformation on the consumer side.

Mainflux [SenML transformer](transformer) is an example of Transformer service for SenML messages.

Mainflux [writers](writers) are using a standalone SenML transformer to preprocess messages before storing them.

[transformers]: https://github.com/mainflux/mainflux/tree/master/transformers/senml
[writers]: https://github.com/mainflux/mainflux/tree/master/writers
