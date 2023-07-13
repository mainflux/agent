// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mainflux/mainflux/pkg/errors"
)

const (
	groupsEndpoint = "groups"
	MaxLevel       = uint64(5)
	MinLevel       = uint64(1)
)

// Group represents the group of Clients.
// Indicates a level in tree hierarchy. Root node is level 1.
// Path in a tree consisting of group IDs
// Paths are unique per owner.
type Group struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id,omitempty"`
	ParentID    string    `json:"parent_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	Level       int       `json:"level,omitempty"`
	Path        string    `json:"path,omitempty"`
	Children    []*Group  `json:"children,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Status      string    `json:"status,omitempty"`
}

func (sdk mfSDK) CreateGroup(g Group, token string) (Group, errors.SDKError) {
	data, err := json.Marshal(g)
	if err != nil {
		return Group{}, errors.NewSDKError(err)
	}
	url := fmt.Sprintf("%s/%s", sdk.usersURL, groupsEndpoint)

	_, body, sdkerr := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), data, http.StatusCreated)
	if sdkerr != nil {
		return Group{}, sdkerr
	}

	g = Group{}
	if err := json.Unmarshal(body, &g); err != nil {
		return Group{}, errors.NewSDKError(err)
	}
	
	return g, nil
}

func (sdk mfSDK) Memberships(clientID string, pm PageMetadata, token string) (MembershipsPage, errors.SDKError) {
	url, err := sdk.withQueryParams(fmt.Sprintf("%s/%s/%s", sdk.usersURL, usersEndpoint, clientID), "memberships", pm)
	if err != nil {
		return MembershipsPage{}, errors.NewSDKError(err)
	}

	_, body, sdkerr := sdk.processRequest(http.MethodGet, url, token, string(CTJSON), nil, http.StatusOK)
	if sdkerr != nil {
		return MembershipsPage{}, sdkerr
	}

	var tp MembershipsPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return MembershipsPage{}, errors.NewSDKError(err)
	}

	return tp, nil
}

func (sdk mfSDK) Groups(pm PageMetadata, token string) (GroupsPage, errors.SDKError) {
	url, err := sdk.withQueryParams(sdk.usersURL, groupsEndpoint, pm)
	if err != nil {
		return GroupsPage{}, errors.NewSDKError(err)
	}
	
	return sdk.getGroups(url, token)
}

func (sdk mfSDK) Parents(id string, pm PageMetadata, token string) (GroupsPage, errors.SDKError) {
	pm.Level = MaxLevel
	url, err := sdk.withQueryParams(fmt.Sprintf("%s/%s/%s", sdk.usersURL, groupsEndpoint, id), "parents", pm)
	if err != nil {
		return GroupsPage{}, errors.NewSDKError(err)
	}
	
	return sdk.getGroups(url, token)
}

func (sdk mfSDK) Children(id string, pm PageMetadata, token string) (GroupsPage, errors.SDKError) {
	pm.Level = MaxLevel
	url, err := sdk.withQueryParams(fmt.Sprintf("%s/%s/%s", sdk.usersURL, groupsEndpoint, id), "children", pm)
	if err != nil {
		return GroupsPage{}, errors.NewSDKError(err)
	}
	
	return sdk.getGroups(url, token)
}

func (sdk mfSDK) getGroups(url, token string) (GroupsPage, errors.SDKError) {
	_, body, err := sdk.processRequest(http.MethodGet, url, token, string(CTJSON), nil, http.StatusOK)
	if err != nil {
		return GroupsPage{}, err
	}

	var tp GroupsPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return GroupsPage{}, errors.NewSDKError(err)
	}
	
	return tp, nil
}

func (sdk mfSDK) Group(id, token string) (Group, errors.SDKError) {
	url := fmt.Sprintf("%s/%s/%s", sdk.usersURL, groupsEndpoint, id)
	_, body, err := sdk.processRequest(http.MethodGet, url, token, string(CTJSON), nil, http.StatusOK)
	if err != nil {
		return Group{}, err
	}

	var t Group
	if err := json.Unmarshal(body, &t); err != nil {
		return Group{}, errors.NewSDKError(err)
	}

	return t, nil
}

func (sdk mfSDK) UpdateGroup(g Group, token string) (Group, errors.SDKError) {
	data, err := json.Marshal(g)
	if err != nil {
		return Group{}, errors.NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s/%s", sdk.usersURL, groupsEndpoint, g.ID)
	_, body, sdkerr := sdk.processRequest(http.MethodPut, url, token, string(CTJSON), data, http.StatusOK)
	if sdkerr != nil {
		return Group{}, sdkerr
	}

	g = Group{}
	if err := json.Unmarshal(body, &g); err != nil {
		return Group{}, errors.NewSDKError(err)
	}

	return g, nil
}

func (sdk mfSDK) EnableGroup(id, token string) (Group, errors.SDKError) {
	return sdk.changeGroupStatus(id, enableEndpoint, token)
}

func (sdk mfSDK) DisableGroup(id, token string) (Group, errors.SDKError) {
	return sdk.changeGroupStatus(id, disableEndpoint, token)
}

func (sdk mfSDK) changeGroupStatus(id, status, token string) (Group, errors.SDKError) {
	url := fmt.Sprintf("%s/%s/%s/%s", sdk.usersURL, groupsEndpoint, id, status)
	_, body, err := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), nil, http.StatusOK)
	if err != nil {
		return Group{}, err
	}
	g := Group{}
	if err := json.Unmarshal(body, &g); err != nil {
		return Group{}, errors.NewSDKError(err)
	}

	return g, nil
}
