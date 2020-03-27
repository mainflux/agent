package util

import "github.com/mainflux/senml"

func EncodeSenML(bn, n, sv string) ([]byte, error) {
	s := senml.Pack{
		Records: []senml.Record{
			senml.Record{
				BaseName:    bn,
				Name:        n,
				StringValue: &sv,
			},
		},
	}
	payload, err := senml.Encode(s, senml.JSON)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
