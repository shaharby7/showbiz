package showbiz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CreateResource creates a new resource in a project.
func (c *Client) CreateResource(ctx context.Context, projectID string, input CreateResourceInput) (*Resource, error) {
	var res Resource
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/projects/%s/resources", projectID), input, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ListResources retrieves a paginated list of resources in a project.
func (c *Client) ListResources(ctx context.Context, projectID string, opts *ListOptions) (*ListResourcesResult, error) {
	path := addListQuery(fmt.Sprintf("/v1/projects/%s/resources", projectID), opts)

	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, path, nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data       []*Resource `json:"data"`
		Pagination pagination  `json:"pagination"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode resources list: %w", err)
	}

	return &ListResourcesResult{
		Data:       envelope.Data,
		NextCursor: envelope.Pagination.NextCursor,
		HasMore:    envelope.Pagination.HasMore,
	}, nil
}

// GetResource retrieves a resource by ID. The resourceID is URL-encoded
// because it may contain colons.
func (c *Client) GetResource(ctx context.Context, projectID, resourceID string) (*Resource, error) {
	var res Resource
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/projects/%s/resources/%s", projectID, url.PathEscape(resourceID)), nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateResource updates a resource.
func (c *Client) UpdateResource(ctx context.Context, projectID, resourceID string, input UpdateResourceInput) (*Resource, error) {
	var res Resource
	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/v1/projects/%s/resources/%s", projectID, url.PathEscape(resourceID)), input, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteResource deletes a resource.
func (c *Client) DeleteResource(ctx context.Context, projectID, resourceID string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/projects/%s/resources/%s", projectID, url.PathEscape(resourceID)), nil, nil)
}
