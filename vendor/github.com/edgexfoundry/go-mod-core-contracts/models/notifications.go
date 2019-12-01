/*******************************************************************************
 * Copyright 2019 Dell Technologies Inc.
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
 *
 *******************************************************************************/

package models

import (
	"encoding/json"
	"reflect"
)

type Notification struct {
	Timestamps
	ID          string                `json:"id,omitempty"`
	Slug        string                `json:"slug,omitempty"`
	Sender      string                `json:"sender,omitempty"`
	Category    NotificationsCategory `json:"category,omitempty"`
	Severity    NotificationsSeverity `json:"severity,omitempty"`
	Content     string                `json:"content,omitempty"`
	Description string                `json:"description,omitempty"`
	Status      NotificationsStatus   `json:"status,omitempty"`
	Labels      []string              `json:"labels,omitempty"`
	ContentType string                `json:"contenttype,omitempty"`
	isValidated bool                  // internal member used for validation check
}

func (n Notification) MarshalJSON() ([]byte, error) {
	test := struct {
		*Timestamps `json:",omitempty"`
		ID          string                `json:"id,omitempty"`
		Slug        string                `json:"slug,omitempty"`
		Sender      string                `json:"sender,omitempty"`
		Category    NotificationsCategory `json:"category,omitempty"`
		Severity    NotificationsSeverity `json:"severity,omitempty"`
		Content     string                `json:"content,omitempty"`
		Description string                `json:"description,omitempty"`
		Status      NotificationsStatus   `json:"status,omitempty"`
		Labels      []string              `json:"labels,omitempty"`
		ContentType string                `json:"contenttype,omitempty"`
	}{
		Timestamps:  &n.Timestamps,
		ID:          n.ID,
		Slug:        n.Slug,
		Sender:      n.Sender,
		Category:    n.Category,
		Severity:    n.Severity,
		Content:     n.Content,
		Description: n.Description,
		Status:      n.Status,
		Labels:      n.Labels,
		ContentType: n.ContentType,
	}

	if reflect.DeepEqual(n.Timestamps, Timestamps{}) {
		test.Timestamps = nil
	}

	return json.Marshal(test)
}

// UnmarshalJSON implements the Unmarshaler interface for the Notification type
func (n *Notification) UnmarshalJSON(data []byte) error {
	var err error
	type Alias struct {
		Timestamps
		ID          *string               `json:"id"`
		Slug        *string               `json:"slug,omitempty,omitempty"`
		Sender      *string               `json:"sender,omitempty"`
		Category    NotificationsCategory `json:"category,omitempty"`
		Severity    NotificationsSeverity `json:"severity,omitempty"`
		Content     *string               `json:"content,omitempty"`
		Description *string               `json:"description,omitempty"`
		Status      NotificationsStatus   `json:"status,omitempty"`
		Labels      []string              `json:"labels,omitempty"`
		ContentType *string               `json:"contenttype,omitempty"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	// Nillable fields
	if a.ID != nil {
		n.ID = *a.ID
	}
	if a.Slug != nil {
		n.Slug = *a.Slug
	}
	if a.Sender != nil {
		n.Sender = *a.Sender
	}
	if a.Content != nil {
		n.Content = *a.Content
	}
	if a.Description != nil {
		n.Description = *a.Description
	}
	if a.ContentType != nil {
		n.ContentType = *a.ContentType
	}
	n.Timestamps = a.Timestamps
	n.Category = a.Category
	n.Severity = a.Severity
	n.Status = a.Status
	n.Labels = a.Labels

	n.isValidated, err = n.Validate()

	return err
}

// Validate satisfies the Validator interface
func (n Notification) Validate() (bool, error) {
	if !n.isValidated {
		if n.ID == "" && n.Slug == "" {
			return false, NewErrContractInvalid("Notifiaction ID and Slug are both blank")
		}
		if n.Sender == "" {
			return false, NewErrContractInvalid("Sender is empty")
		}
		if n.Content == "" {
			return false, NewErrContractInvalid("Content is empty")
		}
		if n.Category == "" {
			return false, NewErrContractInvalid("Category is empty")
		}
		if n.Severity == "" {
			return false, NewErrContractInvalid("Severity is empty")
		}
		if n.Severity != "" && n.Severity != "CRITICAL" && n.Severity != "NORMAL" {
			return false, NewErrContractInvalid("Invalid notification severity")
		}
		if n.Category != "" && n.Category != "SECURITY" && n.Category != "HW_HEALTH" && n.Category != "SW_HEALTH" {
			return false, NewErrContractInvalid("Invalid notification severity")
		}
		if n.Status != "" && n.Status != "NEW" && n.Status != "PROCESSED" && n.Status != "ESCALATED" {
			return false, NewErrContractInvalid("Invalid notification severity")
		}
		err := validate(n)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return n.isValidated, nil
}

/*
 * To String function for Notification Struct
 */
func (n Notification) String() string {
	out, err := json.Marshal(n)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
