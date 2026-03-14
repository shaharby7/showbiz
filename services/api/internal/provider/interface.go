package provider

import "context"

// Provider defines the interface that all infrastructure providers must implement.
type Provider interface {
	Name() string
	ResourceTypes() []string
	ValidateCredentials(ctx context.Context, credentials map[string]interface{}) error
	CreateResource(ctx context.Context, input *CreateResourceInput) (*ResourceOutput, error)
	GetResource(ctx context.Context, resourceID string) (*ResourceOutput, error)
	UpdateResource(ctx context.Context, input *UpdateResourceInput) (*ResourceOutput, error)
	DeleteResource(ctx context.Context, resourceID string) error
	DetectDrifts(ctx context.Context, resources []ResourceExpectedState) ([]DriftReport, error)
}

// CreateResourceInput holds the parameters for creating a resource.
type CreateResourceInput struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
}

// UpdateResourceInput holds the parameters for updating a resource.
type UpdateResourceInput struct {
	ResourceID string                 `json:"resourceId"`
	Properties map[string]interface{} `json:"properties"`
}

// ResourceOutput represents the result of a resource operation.
type ResourceOutput struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	Properties map[string]interface{} `json:"properties"`
}

// ResourceExpectedState describes the expected state of a resource for drift detection.
type ResourceExpectedState struct {
	ResourceID         string                 `json:"resourceId"`
	ExpectedProperties map[string]interface{} `json:"expectedProperties"`
}

// DriftReport describes detected drift for a single resource.
type DriftReport struct {
	ResourceID string                 `json:"resourceId"`
	Drifted    bool                   `json:"drifted"`
	Changes    map[string]interface{} `json:"changes,omitempty"`
}
