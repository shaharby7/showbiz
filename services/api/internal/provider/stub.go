package provider

import (
	"context"

	"github.com/google/uuid"
)

// StubProvider is a mock provider for development and testing.
type StubProvider struct{}

var _ Provider = (*StubProvider)(nil)

func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

func (s *StubProvider) Name() string {
	return "stub"
}

func (s *StubProvider) ResourceTypes() []string {
	return []string{"machine", "network"}
}

func (s *StubProvider) ValidateCredentials(_ context.Context, _ map[string]interface{}) error {
	return nil
}

func (s *StubProvider) CreateResource(_ context.Context, input *CreateResourceInput) (*ResourceOutput, error) {
	return &ResourceOutput{
		ID:         uuid.New().String(),
		Type:       input.Type,
		Name:       input.Name,
		Status:     "running",
		Properties: input.Properties,
	}, nil
}

func (s *StubProvider) GetResource(_ context.Context, resourceID string) (*ResourceOutput, error) {
	return &ResourceOutput{
		ID:         resourceID,
		Type:       "machine",
		Name:       "stub-resource",
		Status:     "running",
		Properties: map[string]interface{}{"cpu": 2, "memory": "4GB"},
	}, nil
}

func (s *StubProvider) UpdateResource(_ context.Context, input *UpdateResourceInput) (*ResourceOutput, error) {
	return &ResourceOutput{
		ID:         input.ResourceID,
		Type:       "machine",
		Name:       "stub-resource",
		Status:     "running",
		Properties: input.Properties,
	}, nil
}

func (s *StubProvider) DeleteResource(_ context.Context, _ string) error {
	return nil
}

func (s *StubProvider) DetectDrifts(_ context.Context, resources []ResourceExpectedState) ([]DriftReport, error) {
	reports := make([]DriftReport, len(resources))
	for i, r := range resources {
		reports[i] = DriftReport{
			ResourceID: r.ResourceID,
			Drifted:    false,
		}
	}
	return reports, nil
}
