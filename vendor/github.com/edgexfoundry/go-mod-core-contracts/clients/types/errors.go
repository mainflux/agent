/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package types

import "fmt"

// ErrNotFound represents an error returned from a service indicating the item being asked for was not found.
type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "Item not found"
}

// ErrResponseNil represents an error returned from a service indicating the response was unexpectedly empty.
type ErrResponseNil struct{}

func (e ErrResponseNil) Error() string {
	return "Response was nil"
}

// ErrServiceClient exposes the details of a service's response in a more granular manner. This is useful when service A
// calls service B and service A needs to make a decision with regard to how it should respond to its own caller based on
// the error thrown from service B.
type ErrServiceClient struct {
	StatusCode int    // StatusCode contains the HTTP status code returned from the target service
	bodyBytes  []byte // bodyBytes contains the response from the target service
	errMsg     string // errMsg contains the error message to be returned. See the Error() method below.
}

// NewErrServiceClient returns an instance of the error interface with ErrServiceClient as its implementation.
func NewErrServiceClient(statusCode int, body []byte) error {
	e := ErrServiceClient{StatusCode: statusCode, bodyBytes: body}
	return e
}

// Error fulfills the error interface and returns an error message assembled from the state of ErrServiceClient.
func (e ErrServiceClient) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.bodyBytes)
}
