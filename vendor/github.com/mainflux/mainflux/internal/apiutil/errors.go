// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package apiutil

import "github.com/mainflux/mainflux/pkg/errors"

// Errors defined in this file are used by the LoggingErrorEncoder decorator
// to distinguish and log API request validation errors and avoid that service
// errors are logged twice.
var (
	// ErrBearerToken indicates missing or invalid bearer user token.
	ErrBearerToken = errors.New("missing or invalid bearer user token")

	// ErrBearerKey indicates missing or invalid bearer entity key.
	ErrBearerKey = errors.New("missing or invalid bearer entity key")

	// ErrMissingID indicates missing entity ID.
	ErrMissingID = errors.New("missing entity id")

	// ErrInvalidAuthKey indicates invalid auth key.
	ErrInvalidAuthKey = errors.New("invalid auth key")

	// ErrInvalidIDFormat indicates an invalid ID format.
	ErrInvalidIDFormat = errors.New("invalid id format provided")

	// ErrNameSize indicates that name size exceeds the max.
	ErrNameSize = errors.New("invalid name size")

	// ErrEmailSize indicates that email size exceeds the max.
	ErrEmailSize = errors.New("invalid email size")

	// ErrInvalidStatus indicates an invalid user account status.
	ErrInvalidStatus = errors.New("invalid user account status")

	// ErrInvalidRole indicates that an invalid role.
	ErrInvalidRole = errors.New("invalid client role")

	// ErrLimitSize indicates that an invalid limit.
	ErrLimitSize = errors.New("invalid limit size")

	// ErrOffsetSize indicates an invalid offset.
	ErrOffsetSize = errors.New("invalid offset size")

	// ErrInvalidOrder indicates an invalid list order.
	ErrInvalidOrder = errors.New("invalid list order provided")

	// ErrInvalidDirection indicates an invalid list direction.
	ErrInvalidDirection = errors.New("invalid list direction provided")

	// ErrEmptyList indicates that entity data is empty.
	ErrEmptyList = errors.New("empty list provided")

	// ErrMalformedPolicy indicates that policies are malformed.
	ErrMalformedPolicy = errors.New("malfalmormed policy")

	// ErrMissingPolicySub indicates that policies are subject.
	ErrMissingPolicySub = errors.New("malfalmormed policy subject")

	// ErrMissingPolicyObj indicates missing policies object.
	ErrMissingPolicyObj = errors.New("malfalmormed policy object")

	// ErrMalformedPolicyAct indicates missing policies action.
	ErrMalformedPolicyAct = errors.New("malfalmormed policy action")

	// ErrMissingCertData indicates missing cert data (ttl).
	ErrMissingCertData = errors.New("missing certificate data")

	// ErrInvalidCertData indicates invalid cert data (ttl).
	ErrInvalidCertData = errors.New("invalid certificate data")

	// ErrInvalidTopic indicates an invalid subscription topic.
	ErrInvalidTopic = errors.New("invalid Subscription topic")

	// ErrInvalidContact indicates an invalid subscription contract.
	ErrInvalidContact = errors.New("invalid Subscription contact")

	// ErrMissingEmail indicates missing email.
	ErrMissingEmail = errors.New("missing email")

	// ErrMissingHost indicates missing host.
	ErrMissingHost = errors.New("missing host")

	// ErrMissingPass indicates missing password.
	ErrMissingPass = errors.New("missing password")

	// ErrMissingConfPass indicates missing conf password.
	ErrMissingConfPass = errors.New("missing conf password")

	// ErrInvalidResetPass indicates an invalid reset password.
	ErrInvalidResetPass = errors.New("invalid reset password")

	// ErrInvalidComparator indicates an invalid comparator.
	ErrInvalidComparator = errors.New("invalid comparator")

	// ErrMissingMemberType indicates missing group member type.
	ErrMissingMemberType = errors.New("missing group member type")

	// ErrInvalidAPIKey indicates an invalid API key type.
	ErrInvalidAPIKey = errors.New("invalid api key type")

	// ErrMaxLevelExceeded indicates an invalid group level.
	ErrMaxLevelExceeded = errors.New("invalid group level (should be lower than 5)")

	// ErrBootstrapState indicates an invalid boostrap state.
	ErrBootstrapState = errors.New("invalid bootstrap state")

	// ErrMissingIdentity indicates missing entity Identity.
	ErrMissingIdentity = errors.New("missing entity identity")

	// ErrMissingSecret indicates missing secret.
	ErrMissingSecret = errors.New("missing secret")

	// ErrMissingOwner indicates missing entity owner.
	ErrMissingOwner = errors.New("missing entity owner")

	// ErrMissingPolicyOwner indicated malformed policy owner.
	ErrMissingPolicyOwner = errors.New("malformed policy owner")

	// ErrMissingPolicyEntityType indicates malformed policy entity type.
	ErrMissingPolicyEntityType = errors.New("malformed or missing entity type")

	// ErrMissingName indicates missing identity name.
	ErrMissingName = errors.New("missing identity name")

	// ErrInvalidSecret indicates invalid secret.
	ErrInvalidSecret = errors.New("missing secret")

	// ErrInvalidLevel indicates an invalid group level.
	ErrInvalidLevel = errors.New("invalid group level (should be between 0 and 5)")

	// ErrNotFoundParam indicates that the parameter was not found in the query.
	ErrNotFoundParam = errors.New("parameter not found in the query")

	// ErrMalformedEntity indicates a malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrInvalidQueryParams indicates invalid query parameters.
	ErrInvalidQueryParams = errors.New("invalid query parameters")

	// ErrInvalidVisibilityType indicates invalid visibility type.
	ErrInvalidVisibilityType = errors.New("invalid visibility type")
)
