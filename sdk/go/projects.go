package showbiz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateProject creates a new project in an organization.
func (c *Client) CreateProject(ctx context.Context, orgID string, input CreateProjectInput) (*Project, error) {
	var project Project
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/projects", orgID), input, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// ListProjects retrieves a paginated list of projects in an organization.
func (c *Client) ListProjects(ctx context.Context, orgID string, opts *ListOptions) (*ListProjectsResult, error) {
	path := addListQuery(fmt.Sprintf("/v1/organizations/%s/projects", orgID), opts)

	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, path, nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data       []*Project `json:"data"`
		Pagination pagination `json:"pagination"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode projects list: %w", err)
	}

	return &ListProjectsResult{
		Data:       envelope.Data,
		NextCursor: envelope.Pagination.NextCursor,
		HasMore:    envelope.Pagination.HasMore,
	}, nil
}

// GetProject retrieves a project by ID.
func (c *Client) GetProject(ctx context.Context, orgID, projectID string) (*Project, error) {
	var project Project
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s/projects/%s", orgID, projectID), nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates a project.
func (c *Client) UpdateProject(ctx context.Context, orgID, projectID string, input UpdateProjectInput) (*Project, error) {
	var project Project
	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/v1/organizations/%s/projects/%s", orgID, projectID), input, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// DeleteProject deletes a project.
func (c *Client) DeleteProject(ctx context.Context, orgID, projectID string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/projects/%s", orgID, projectID), nil, nil)
}
