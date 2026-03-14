package showbiz

import (
	"context"
	"fmt"
	"net/http"
)

// ListProviders retrieves all available providers.
func (c *Client) ListProviders(ctx context.Context) ([]*Provider, error) {
	var providers []*Provider
	if err := c.do(ctx, http.MethodGet, "/v1/providers", nil, &providers); err != nil {
		return nil, err
	}
	return providers, nil
}

// GetProvider retrieves a provider by ID (name).
func (c *Client) GetProvider(ctx context.Context, id string) (*Provider, error) {
	var provider Provider
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/providers/%s", id), nil, &provider); err != nil {
		return nil, err
	}
	return &provider, nil
}
