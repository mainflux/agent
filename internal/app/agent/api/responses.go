// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

type genericRes struct {
	Service  string `json:"service"`
	Response string `json:"response"`
}

type execRes struct {
	BaseName string `json:"bn"`
	Name     string `json:"n"`
	Value    string `json:"vs"`
}
