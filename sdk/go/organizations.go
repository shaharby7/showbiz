package showbiz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateOrganization creates a new organization.
func (c *Client) CreateOrganization(ctx context.Context, input CreateOrganizationInput) (*Organization, error) {
	var org Organization
	if err := c.do(ctx, http.MethodPost, "/v1/organizations", input, &org); err != nil {
		return nil, err
	}
	return &org, nil
}

// GetOrganization retrieves an organization by ID.
func (c *Client) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	var org Organization
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s", id), nil, &org); err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrganizations retrieves a paginated list of organizations.
func (c *Client) ListOrganizations(ctx context.Context, opts *ListOptions) (*ListOrganizationsResult, error) {
	path := addListQuery("/v1/organizations", opts)

	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, path, nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data       []*Organization `json:"data"`
		Pagination pagination      `json:"pagination"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode organizations list: %w", err)
	}

	return &ListOrganizationsResult{
		Data:       envelope.Data,
		NextCursor: envelope.Pagination.NextCursor,
		HasMore:    envelope.Pagination.HasMore,
	}, nil
}

// UpdateOrganization updates an organization.
func (c *Client) UpdateOrganization(ctx context.Context, id string, input UpdateOrganizationInput) (*Organization, error) {
	var org Organization
	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/v1/organizations/%s", id), input, &org); err != nil {
		return nil, err
	}
	return &org, nil
}

// DeactivateOrganization deactivates an organization.
func (c *Client) DeactivateOrganization(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/deactivate", id), nil, nil)
}

// ActivateOrganization activates an organization.
func (c *Client) ActivateOrganization(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/activate", id), nil, nil)
}

// ListMembers returns the members of an organization.
func (c *Client) ListMembers(ctx context.Context, orgID string) ([]*User, error) {
	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s/members", orgID), nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data []*User `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode members list: %w", err)
	}

	return envelope.Data, nil
}

// AddMember adds a member to an organization by email.
func (c *Client) AddMember(ctx context.Context, orgID string, email string) error {
	payload := struct {
		Email string `json:"email"`
	}{Email: email}
	return c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/members", orgID), payload, nil)
}

// RemoveMember removes a member from an organization by email.
func (c *Client) RemoveMember(ctx context.Context, orgID string, email string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/members/%s", orgID, email), nil, nil)
}
