package service

import (
	"context"
	"fmt"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
	"github.com/showbiz-io/showbiz/services/api/internal/provider"
	"github.com/showbiz-io/showbiz/services/api/internal/repository"
)

type ResourceService struct {
	resourceRepo *repository.ResourceRepo
	connRepo     *repository.ConnectionRepo
	registry     *provider.Registry
}

func NewResourceService(resourceRepo *repository.ResourceRepo, connRepo *repository.ConnectionRepo, registry *provider.Registry) *ResourceService {
	return &ResourceService{
		resourceRepo: resourceRepo,
		connRepo:     connRepo,
		registry:     registry,
	}
}

type CreateResourceInput struct {
	Name         string                 `json:"name"`
	ConnectionID string                 `json:"connectionId"`
	ResourceType string                 `json:"resourceType"`
	Values       map[string]interface{} `json:"values"`
}

type UpdateResourceInput struct {
	Values map[string]interface{} `json:"values"`
}

func (s *ResourceService) Create(ctx context.Context, projectID string, input CreateResourceInput) (*model.Resource, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.ConnectionID == "" {
		return nil, fmt.Errorf("connectionId is required")
	}
	if input.ResourceType == "" {
		return nil, fmt.Errorf("resourceType is required")
	}

	// Validate connection exists
	conn, err := s.connRepo.GetConnectionByID(ctx, input.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("get connection: %w", err)
	}
	if conn == nil {
		return nil, fmt.Errorf("connection not found")
	}

	// Validate resource type against provider
	p, ok := s.registry.Get(conn.Provider)
	if !ok {
		return nil, fmt.Errorf("provider not found")
	}
	validType := false
	for _, rt := range p.ResourceTypes() {
		if rt == input.ResourceType {
			validType = true
			break
		}
	}
	if !validType {
		return nil, fmt.Errorf("invalid resource type")
	}

	// Check name uniqueness within the project
	existing, err := s.resourceRepo.GetResourceByProjectAndName(ctx, projectID, input.Name)
	if err != nil {
		return nil, fmt.Errorf("check name uniqueness: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("resource name already exists")
	}

	// Generate deterministic ID: sbz:<resourceType>:<projectId>:<connectionName>:<resourceName>
	id := fmt.Sprintf("sbz:%s:%s:%s:%s", input.ResourceType, projectID, conn.Name, input.Name)

	res := &model.Resource{
		ID:           id,
		Name:         input.Name,
		ProjectID:    projectID,
		ConnectionID: input.ConnectionID,
		ResourceType: input.ResourceType,
		Values:       input.Values,
		Status:       "creating",
	}

	if err := s.resourceRepo.CreateResource(ctx, res); err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	// No actual provider call yet — immediately set to active
	if err := s.resourceRepo.UpdateResource(ctx, id, input.Values, "active"); err != nil {
		return nil, fmt.Errorf("activate resource: %w", err)
	}

	created, err := s.resourceRepo.GetResourceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch created resource: %w", err)
	}
	return created, nil
}

func (s *ResourceService) Get(ctx context.Context, id string) (*model.Resource, error) {
	res, err := s.resourceRepo.GetResourceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}
	if res == nil {
		return nil, fmt.Errorf("resource not found")
	}
	return res, nil
}

func (s *ResourceService) List(ctx context.Context, projectID string, cursor string, limit int) ([]*model.Resource, string, error) {
	resources, nextCursor, err := s.resourceRepo.ListResourcesByProjectID(ctx, projectID, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("list resources: %w", err)
	}
	return resources, nextCursor, nil
}

func (s *ResourceService) Update(ctx context.Context, id string, input UpdateResourceInput) (*model.Resource, error) {
	if input.Values == nil {
		return nil, fmt.Errorf("values is required")
	}

	res, err := s.resourceRepo.GetResourceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}
	if res == nil {
		return nil, fmt.Errorf("resource not found")
	}

	if err := s.resourceRepo.UpdateResource(ctx, id, input.Values, res.Status); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	updated, err := s.resourceRepo.GetResourceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch updated resource: %w", err)
	}
	return updated, nil
}

func (s *ResourceService) Delete(ctx context.Context, id string) error {
	res, err := s.resourceRepo.GetResourceByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get resource: %w", err)
	}
	if res == nil {
		return fmt.Errorf("resource not found")
	}

	if err := s.resourceRepo.DeleteResource(ctx, id); err != nil {
		return fmt.Errorf("delete resource: %w", err)
	}
	return nil
}
