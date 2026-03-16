package model

import "time"

type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type User struct {
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	OrganizationID string    `json:"organizationId"`
	DisplayName    string    `json:"displayName"`
	EmailVerified  bool      `json:"emailVerified"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Project struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organizationId"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Connection struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	ProjectID   string                 `json:"projectId"`
	Provider    string                 `json:"provider"`
	Credentials map[string]interface{} `json:"credentials,omitempty"` // write-only, omitted in responses
	Config      map[string]interface{} `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type Resource struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	ProjectID    string                 `json:"projectId"`
	ConnectionID *string                `json:"connectionId"` // nil for Showbiz-managed resource types
	ResourceType string                 `json:"resourceType"`
	Values       map[string]interface{} `json:"values"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type Policy struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Scope          string    `json:"scope"`
	OrganizationID string    `json:"organizationId,omitempty"`
	Permissions    []string  `json:"permissions"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type PolicyAttachment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	UserEmail string    `json:"userEmail"`
	PolicyID  string    `json:"policyId"`
	CreatedAt time.Time `json:"createdAt"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// Pagination represents cursor-based pagination.
type Pagination struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// ListResponse wraps a list of items with pagination.
type ListResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
