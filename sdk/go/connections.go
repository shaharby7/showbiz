package showbiz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateConnection creates a new connection in a project.
func (c *Client) CreateConnection(ctx context.Context, projectID string, input CreateConnectionInput) (*Connection, error) {
	var conn Connection
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/projects/%s/connections", projectID), input, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

// ListConnections retrieves a paginated list of connections in a project.
func (c *Client) ListConnections(ctx context.Context, projectID string, opts *ListOptions) (*ListConnectionsResult, error) {
	path := addListQuery(fmt.Sprintf("/v1/projects/%s/connections", projectID), opts)

	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, path, nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data       []*Connection `json:"data"`
		Pagination pagination    `json:"pagination"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode connections list: %w", err)
	}

	return &ListConnectionsResult{
		Data:       envelope.Data,
		NextCursor: envelope.Pagination.NextCursor,
		HasMore:    envelope.Pagination.HasMore,
	}, nil
}

// GetConnection retrieves a connection by ID.
func (c *Client) GetConnection(ctx context.Context, projectID, connectionID string) (*Connection, error) {
	var conn Connection
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/projects/%s/connections/%s", projectID, connectionID), nil, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

// UpdateConnection updates a connection.
func (c *Client) UpdateConnection(ctx context.Context, projectID, connectionID string, input UpdateConnectionInput) (*Connection, error) {
	var conn Connection
	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/v1/projects/%s/connections/%s", projectID, connectionID), input, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

// DeleteConnection deletes a connection.
func (c *Client) DeleteConnection(ctx context.Context, projectID, connectionID string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/projects/%s/connections/%s", projectID, connectionID), nil, nil)
}
