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

package clients

// Do not assume that if a constant is identified by your IDE as not being used within this module that it is not being
// used at all. Any application wishing to exchange information with the EdgeX core services will utilize this module,
// so constants located here may be used externally.
//
// Miscellaneous constants
const (
	ClientMonitorDefault = 15000            // Defaults the interval at which a given service client will refresh its endpoint from the Registry, if used
	CorrelationHeader    = "correlation-id" // Sets the key of the Correlation ID HTTP header
)

// Constants related to defined routes in the service APIs
const (
	ApiVersionRoute            = "/api/version"
	ApiBase                    = "/api/v1"
	ApiAddressableRoute        = "/api/v1/addressable"
	ApiCallbackRoute           = "/api/v1/callback"
	ApiCommandRoute            = "/api/v1/command"
	ApiConfigRoute             = "/api/v1/config"
	ApiDeviceRoute             = "/api/v1/device"
	ApiDeviceProfileRoute      = "/api/v1/deviceprofile"
	ApiDeviceServiceRoute      = "/api/v1/deviceservice"
	ApiEventRoute              = "/api/v1/event"
	ApiHealthRoute             = "/api/v1/health"
	ApiLoggingRoute            = "/api/v1/logs"
	ApiMetricsRoute            = "/api/v1/metrics"
	ApiNotificationRoute       = "/api/v1/notification"
	ApiNotifyRegistrationRoute = "/api/v1/notify/registrations"
	ApiOperationRoute          = "/api/v1/operation"
	ApiPingRoute               = "/api/v1/ping"
	ApiProvisionWatcherRoute   = "/api/v1/provisionwatcher"
	ApiReadingRoute            = "/api/v1/reading"
	ApiRegistrationRoute       = "/api/v1/registration"
	ApiRegistrationByNameRoute = ApiRegistrationRoute + "/name"
	ApiSubscriptionRoute       = "/api/v1/subscription"
	ApiTransmissionRoute       = "/api/v1/transmission"
	ApiValueDescriptorRoute    = "/api/v1/valuedescriptor"
	ApiIntervalRoute           = "/api/v1/interval"
	ApiIntervalActionRoute     = "/api/v1/intervalaction"
)

// Constants related to how services identify themselves in the Service Registry
const (
	ServiceKeyPrefix                    = "edgex-"
	ConfigSeedServiceKey                = "edgex-config-seed"
	CoreCommandServiceKey               = "edgex-core-command"
	CoreDataServiceKey                  = "edgex-core-data"
	CoreMetaDataServiceKey              = "edgex-core-metadata"
	SupportLoggingServiceKey            = "edgex-support-logging"
	SupportNotificationsServiceKey      = "edgex-support-notifications"
	SystemManagementAgentServiceKey     = "edgex-sys-mgmt-agent"
	SupportSchedulerServiceKey          = "edgex-support-scheduler"
	SecuritySecretStoreSetupServiceKey  = "edgex-security-secretstore-setup"
	SecuritySecretsSetupServiceKey      = "edgex-security-secrets-setup"
	SecurityProxySetupServiceKey        = "edgex-security-proxy-setup"
	SecurityFileTokenProviderServiceKey = "edgex-security-file-token-provider"
)

// Constants related to the possible content types supported by the APIs
const (
	ContentType     = "Content-Type"
	ContentTypeCBOR = "application/cbor"
	ContentTypeJSON = "application/json"
	ContentTypeYAML = "application/x-yaml"
	ContentTypeText = "text/plain"
)
