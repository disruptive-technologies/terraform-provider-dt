// Copyright (c) HashiCorp, Inc.

package dt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateContactGroupRequest struct {
	Organization string       `json:"organization"`
	ContactGroup ContactGroup `json:"contactGroup"`
}

type UpdateContactGroupRequest struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type ContactGroup struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	ContactCount int32  `json:"contactCount"`
}

func (c *Client) GetContactGroup(name string) (ContactGroup, error) {
	organizationID, scheduledExportID, err := ParseResourceName(name)
	if err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to parse resource name: %w", err)
	}
	url := fmt.Sprintf("%s/v2/organizations/%s/contactGroups/%s", strings.TrimSuffix(c.URL, "/"), organizationID, scheduledExportID)

	responseBody, err := c.DoRequest(context.Background(), http.MethodGet, url, nil, nil)
	if err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to get contact group: %w", err)
	}

	var contactGroup ContactGroup
	if err := json.Unmarshal(responseBody, &contactGroup); err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to unmarshal contact group: %w", err)
	}

	return contactGroup, nil
}

func (c *Client) CreateContactGroup(ctx context.Context, request CreateContactGroupRequest) (ContactGroup, error) {
	url := fmt.Sprintf("%s/v2/%s/contactGroups", strings.TrimSuffix(c.URL, "/"), request.Organization)

	body, err := json.Marshal(request.ContactGroup)
	if err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to marshal contact group: %w", err)
	}

	responseBody, err := c.DoRequest(ctx, http.MethodPost, url, body, nil)
	if err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to create contact group: %w", err)
	}

	var createdGroup ContactGroup
	if err := json.Unmarshal(responseBody, &createdGroup); err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to unmarshal created contact group: %w", err)
	}

	return createdGroup, nil
}

func (c *Client) UpdateContactGroup(ctx context.Context, group UpdateContactGroupRequest, name string) (ContactGroup, error) {
	url := fmt.Sprintf("%s/v2/%s", strings.TrimSuffix(c.URL, "/"), name)

	body, err := json.Marshal(group)
	if err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to marshal contact group: %w", err)
	}

	responseBody, err := c.DoRequest(ctx, http.MethodPatch, url, body, nil)
	if err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to update contact group: %w", err)
	}

	var updatedGroup ContactGroup
	if err := json.Unmarshal(responseBody, &updatedGroup); err != nil {
		return ContactGroup{}, fmt.Errorf("dt: failed to unmarshal updated contact group: %w", err)
	}

	return updatedGroup, nil
}

func (c *Client) DeleteContactGroup(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/v2/%s", strings.TrimSuffix(c.URL, "/"), name)

	_, err := c.DoRequest(ctx, http.MethodDelete, url, nil, nil)
	if err != nil {
		return fmt.Errorf("dt: failed to delete contact group: %w", err)
	}

	return nil
}
