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

package models

// ErrContractInvalid is a specific error type for handling model validation failures. Type checking within
// the calling application will facilitate more explicit error handling whereby it's clear that validation
// has failed as opposed to something unexpected happening.
type ErrContractInvalid struct {
	errMsg string
}

// NewErrContractInvalid returns an instance of the error interface with ErrContractInvalid as its implementation.
func NewErrContractInvalid(message string) error {
	return ErrContractInvalid{errMsg: message}
}

// Error fulfills the error interface and returns an error message assembled from the state of ErrContractInvalid.
func (e ErrContractInvalid) Error() string {
	return e.errMsg
}
