// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package apiutil

// ErrorRes represents the HTTP error response body.
type ErrorRes struct {
	Err string `json:"error"`
}
