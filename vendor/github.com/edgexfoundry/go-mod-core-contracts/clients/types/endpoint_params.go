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

/*
 Package types provides supporting types that facilitate the various service client implementations.
*/
package types

// EndpointParams is a type that allows for the passing of common parameters to service clients
// for initialization.
type EndpointParams struct {
	ServiceKey  string // The key of the service as found in the service registry (e.g. Consul)
	Path        string // The path to the service's endpoint following port number in the URL
	UseRegistry bool   // An indication of whether or not endpoint information should be obtained from a service registry provider.
	Url         string // If a service registry is not being used, then provide the full URL endpoint
	Interval    int    // The interval in milliseconds governing how often the client polls to keep the endpoint current
}
