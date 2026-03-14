package showbiz

import (
	"context"
	"net/http"
)

// Register creates a new user account.
func (c *Client) Register(ctx context.Context, input RegisterInput) (*User, error) {
	var user User
	if err := c.do(ctx, http.MethodPost, "/v1/auth/register", input, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Login authenticates a user and stores the resulting tokens on the client.
func (c *Client) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	var resp AuthResponse
	if err := c.do(ctx, http.MethodPost, "/v1/auth/login", input, &resp); err != nil {
		return nil, err
	}
	c.setTokens(resp.AccessToken, resp.RefreshToken)
	return &resp, nil
}

// RefreshToken exchanges the stored refresh token for new tokens.
func (c *Client) RefreshToken(ctx context.Context) (*AuthResponse, error) {
	rt := c.getRefreshToken()
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{RefreshToken: rt}

	var resp AuthResponse
	if err := c.do(ctx, http.MethodPost, "/v1/auth/refresh", payload, &resp); err != nil {
		return nil, err
	}
	c.setTokens(resp.AccessToken, resp.RefreshToken)
	return &resp, nil
}

// Me returns the currently authenticated user.
func (c *Client) Me(ctx context.Context) (*User, error) {
	var user User
	if err := c.do(ctx, http.MethodGet, "/v1/auth/me", nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
