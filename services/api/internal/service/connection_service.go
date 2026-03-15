package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaharby7/showbiz/services/api/internal/model"
	"github.com/shaharby7/showbiz/services/api/internal/provider"
	"github.com/shaharby7/showbiz/services/api/internal/repository"
)

type ConnectionService struct {
	connRepo *repository.ConnectionRepo
	registry *provider.Registry
}

func NewConnectionService(connRepo *repository.ConnectionRepo, registry *provider.Registry) *ConnectionService {
	return &ConnectionService{
		connRepo: connRepo,
		registry: registry,
	}
}

type CreateConnectionInput struct {
	Name        string                 `json:"name"`
	Provider    string                 `json:"provider"`
	Credentials map[string]interface{} `json:"credentials"`
	Config      map[string]interface{} `json:"config"`
}

type UpdateConnectionInput struct {
	Config map[string]interface{} `json:"config"`
}

func (s *ConnectionService) Create(ctx context.Context, projectID string, input CreateConnectionInput) (*model.Connection, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}

	if _, ok := s.registry.Get(input.Provider); !ok {
		return nil, fmt.Errorf("provider not found")
	}

	conn := &model.Connection{
		ID:          uuid.New().String(),
		Name:        input.Name,
		ProjectID:   projectID,
		Provider:    input.Provider,
		Credentials: input.Credentials,
		Config:      input.Config,
	}

	if err := s.connRepo.CreateConnection(ctx, conn); err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	created, err := s.connRepo.GetConnectionByID(ctx, conn.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created connection: %w", err)
	}
	return created, nil
}

func (s *ConnectionService) Get(ctx context.Context, id string) (*model.Connection, error) {
	conn, err := s.connRepo.GetConnectionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get connection: %w", err)
	}
	if conn == nil {
		return nil, fmt.Errorf("connection not found")
	}
	return conn, nil
}

func (s *ConnectionService) List(ctx context.Context, projectID string, cursor string, limit int) ([]*model.Connection, string, error) {
	conns, nextCursor, err := s.connRepo.ListConnectionsByProjectID(ctx, projectID, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("list connections: %w", err)
	}
	return conns, nextCursor, nil
}

func (s *ConnectionService) Update(ctx context.Context, id string, input UpdateConnectionInput) (*model.Connection, error) {
	if input.Config == nil {
		return nil, fmt.Errorf("config is required")
	}

	if err := s.connRepo.UpdateConnection(ctx, id, input.Config); err != nil {
		return nil, fmt.Errorf("update connection: %w", err)
	}

	updated, err := s.connRepo.GetConnectionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch updated connection: %w", err)
	}
	return updated, nil
}

func (s *ConnectionService) Delete(ctx context.Context, id string) error {
	conn, err := s.connRepo.GetConnectionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get connection: %w", err)
	}
	if conn == nil {
		return fmt.Errorf("connection not found")
	}

	hasResources, err := s.connRepo.HasResources(ctx, id)
	if err != nil {
		return fmt.Errorf("check resources: %w", err)
	}
	if hasResources {
		return fmt.Errorf("connection has resources")
	}

	if err := s.connRepo.DeleteConnection(ctx, id); err != nil {
		return fmt.Errorf("delete connection: %w", err)
	}
	return nil
}
