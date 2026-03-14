package showbiz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Client is the Showbiz API client.
type Client struct {
	baseURL    string
	httpClient *http.Client

	mu           sync.RWMutex
	token        string
	refreshToken string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithToken sets a pre-existing bearer token on the client.
func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

// NewClient creates a new Showbiz API client.
func NewClient(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// setTokens stores the access and refresh tokens on the client.
func (c *Client) setTokens(access, refresh string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = access
	c.refreshToken = refresh
}

// getToken returns the current access token.
func (c *Client) getToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

// getRefreshToken returns the current refresh token.
func (c *Client) getRefreshToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.refreshToken
}

// do executes an HTTP request against the Showbiz API.
// If result is non-nil, the response body is decoded into it.
// On 401 responses with a stored refresh token, it attempts a single token refresh and retries.
func (c *Client) do(ctx context.Context, method, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	// Handle 401 with automatic token refresh.
	if resp.StatusCode == http.StatusUnauthorized && c.getRefreshToken() != "" {
		resp.Body.Close()

		if err := c.doRefresh(ctx); err != nil {
			return err
		}

		resp, err = c.doRequest(ctx, method, path, body)
		if err != nil {
			return err
		}
	}

	// Check for error responses.
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp)
	}

	// 204 No Content — nothing to decode.
	if resp.StatusCode == http.StatusNoContent {
		resp.Body.Close()
		return nil
	}

	if result != nil {
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(result)
	}

	resp.Body.Close()
	return nil
}

// doRequest builds and executes a single HTTP request (no retry logic).
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("showbiz: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("showbiz: failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token := c.getToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return c.httpClient.Do(req)
}

// doRefresh performs a token refresh using the stored refresh token.
func (c *Client) doRefresh(ctx context.Context) error {
	rt := c.getRefreshToken()
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{RefreshToken: rt}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v1/auth/refresh", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp)
	}

	defer resp.Body.Close()
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("showbiz: failed to decode refresh response: %w", err)
	}

	c.setTokens(authResp.AccessToken, authResp.RefreshToken)
	return nil
}

// addListQuery appends cursor and limit query parameters to a path.
func addListQuery(path string, opts *ListOptions) string {
	if opts == nil {
		return path
	}
	v := url.Values{}
	if opts.Cursor != "" {
		v.Set("cursor", opts.Cursor)
	}
	if opts.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if len(v) == 0 {
		return path
	}
	return path + "?" + v.Encode()
}
