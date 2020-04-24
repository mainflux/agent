/*******************************************************************************
 * Copyright 2020 Dell Inc.
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

package interfaces

// URLClient is the interface for types that need to define some way to retrieve URLClient information about services.
// This information can be anything that must be determined at runtime, whether it is unknown or simply not yet known.
type URLClient interface {
	// Prefix returns the URLClient base path (or root) of a service.
	// This is the common root of all REST calls to the service,
	// and is defined on a per service (rather than per endpoint) basis.
	// Prefix returns the root URLClient for REST calls to the service if it was able to retrieve that URLClient;
	// it returns an error otherwise.
	Prefix() (string, error)
}
