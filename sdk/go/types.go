package showbiz

import "time"

// Organization represents a Showbiz organization.
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// User represents a Showbiz user.
type User struct {
	Email          string    `json:"email"`
	OrganizationID string    `json:"organizationId"`
	DisplayName    string    `json:"displayName"`
	EmailVerified  bool      `json:"emailVerified"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Project represents a project within an organization.
type Project struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organizationId"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Connection represents a provider connection within a project.
type Connection struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	ProjectID   string                 `json:"projectId"`
	Provider    string                 `json:"provider"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// Resource represents a managed resource within a project.
type Resource struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	ProjectID    string                 `json:"projectId"`
	ConnectionID string                 `json:"connectionId"`
	ResourceType string                 `json:"resourceType"`
	Values       map[string]interface{} `json:"values"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// Policy represents an IAM policy.
type Policy struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Scope          string    `json:"scope"`
	OrganizationID string    `json:"organizationId,omitempty"`
	Permissions    []string  `json:"permissions"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// PolicyAttachment represents a policy attached to a user in a project.
type PolicyAttachment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	UserEmail string    `json:"userEmail"`
	PolicyID  string    `json:"policyId"`
	CreatedAt time.Time `json:"createdAt"`
}

// Provider represents a resource provider.
type Provider struct {
	Name          string   `json:"name"`
	ResourceTypes []string `json:"resourceTypes"`
}

// AuthResponse is the response from login and refresh endpoints.
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ListOptions configures pagination for list endpoints.
type ListOptions struct {
	Cursor string
	Limit  int
}

// pagination is the API pagination envelope.
type pagination struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// Concrete ListResult types for paginated responses.

// ListOrganizationsResult is a paginated list of organizations.
type ListOrganizationsResult struct {
	Data       []*Organization `json:"data"`
	NextCursor string
	HasMore    bool
}

// ListProjectsResult is a paginated list of projects.
type ListProjectsResult struct {
	Data       []*Project `json:"data"`
	NextCursor string
	HasMore    bool
}

// ListConnectionsResult is a paginated list of connections.
type ListConnectionsResult struct {
	Data       []*Connection `json:"data"`
	NextCursor string
	HasMore    bool
}

// ListResourcesResult is a paginated list of resources.
type ListResourcesResult struct {
	Data       []*Resource `json:"data"`
	NextCursor string
	HasMore    bool
}

// Input types for create and update operations.

// RegisterInput is the input for user registration.
type RegisterInput struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	DisplayName    string `json:"displayName,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
}

// LoginInput is the input for user login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateOrganizationInput is the input for creating an organization.
type CreateOrganizationInput struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// UpdateOrganizationInput is the input for updating an organization.
type UpdateOrganizationInput struct {
	DisplayName string `json:"displayName"`
}

// CreateProjectInput is the input for creating a project.
type CreateProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateProjectInput is the input for updating a project.
type UpdateProjectInput struct {
	Description string `json:"description"`
}

// CreateConnectionInput is the input for creating a connection.
type CreateConnectionInput struct {
	Name        string                 `json:"name"`
	Provider    string                 `json:"provider"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// UpdateConnectionInput is the input for updating a connection.
type UpdateConnectionInput struct {
	Config map[string]interface{} `json:"config"`
}

// CreateResourceInput is the input for creating a resource.
type CreateResourceInput struct {
	Name         string                 `json:"name"`
	ConnectionID string                 `json:"connectionId"`
	ResourceType string                 `json:"resourceType"`
	Values       map[string]interface{} `json:"values,omitempty"`
}

// UpdateResourceInput is the input for updating a resource.
type UpdateResourceInput struct {
	Values map[string]interface{} `json:"values"`
}

// CreatePolicyInput is the input for creating an organization policy.
type CreatePolicyInput struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// AttachPolicyInput is the input for attaching a policy to a user in a project.
type AttachPolicyInput struct {
	UserEmail string `json:"userEmail"`
	PolicyID  string `json:"policyId"`
}

// DetachPolicyInput is the input for detaching a policy from a user in a project.
type DetachPolicyInput struct {
	UserEmail string `json:"userEmail"`
	PolicyID  string `json:"policyId"`
}
