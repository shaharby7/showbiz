package service

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/shaharby7/showbiz/services/api/internal/model"
	"github.com/shaharby7/showbiz/services/api/internal/provider"
	"github.com/shaharby7/showbiz/services/api/internal/repository"
	"github.com/shaharby7/showbiz/services/api/internal/resource"
)

type ResourceService struct {
	resourceRepo     *repository.ResourceRepo
	connRepo         *repository.ConnectionRepo
	registry         *provider.Registry
	resourceTypeReg  *resource.Registry
}

func NewResourceService(resourceRepo *repository.ResourceRepo, connRepo *repository.ConnectionRepo, registry *provider.Registry, resourceTypeReg *resource.Registry) *ResourceService {
	return &ResourceService{
		resourceRepo:    resourceRepo,
		connRepo:        connRepo,
		registry:        registry,
		resourceTypeReg: resourceTypeReg,
	}
}

type CreateResourceInput struct {
	Name         string                 `json:"name"`
	ConnectionID string                 `json:"connectionId,omitempty"`
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
	if input.ResourceType == "" {
		return nil, fmt.Errorf("resourceType is required")
	}

	// Look up resource type
	rt, ok := s.resourceTypeReg.Get(input.ResourceType)
	if !ok {
		return nil, fmt.Errorf("unknown resource type")
	}

	// Validate input values against the resource type schema
	if err := rt.ValidateCreate(input.Values); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check name uniqueness within the project
	existing, err := s.resourceRepo.GetResourceByProjectAndName(ctx, projectID, input.Name)
	if err != nil {
		return nil, fmt.Errorf("check name uniqueness: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("resource name already exists")
	}

	var connID *string
	var id string

	if rt.RequiresConnection() {
		// Provider-backed resource: connection is required
		if input.ConnectionID == "" {
			return nil, fmt.Errorf("connectionId is required for resource type %q", input.ResourceType)
		}

		conn, err := s.connRepo.GetConnectionByID(ctx, input.ConnectionID)
		if err != nil {
			return nil, fmt.Errorf("get connection: %w", err)
		}
		if conn == nil {
			return nil, fmt.Errorf("connection not found")
		}

		// Validate provider supports this resource type
		p, ok := s.registry.Get(conn.Provider)
		if !ok {
			return nil, fmt.Errorf("provider not found")
		}
		validType := false
		for _, prt := range p.ResourceTypes() {
			if prt == input.ResourceType {
				validType = true
				break
			}
		}
		if !validType {
			return nil, fmt.Errorf("invalid resource type")
		}

		connID = &input.ConnectionID
		id = fmt.Sprintf("sbz:%s:%s:%s:%s", input.ResourceType, projectID, conn.Name, input.Name)

		// Call provider to create the resource
		providerOutput, err := p.CreateResource(ctx, &provider.CreateResourceInput{
			Type:       input.ResourceType,
			Name:       input.Name,
			Properties: input.Values,
		})
		if err != nil {
			return nil, fmt.Errorf("provider create: %w", err)
		}

		// Merge provider output into values
		values := input.Values
		if values == nil {
			values = make(map[string]interface{})
		}
		if providerOutput.Properties != nil {
			for k, v := range providerOutput.Properties {
				values[k] = v
			}
		}

		res := &model.Resource{
			ID:           id,
			Name:         input.Name,
			ProjectID:    projectID,
			ConnectionID: connID,
			ResourceType: input.ResourceType,
			Values:       values,
			Status:       providerOutput.Status,
		}

		if err := s.resourceRepo.CreateResource(ctx, res); err != nil {
			return nil, fmt.Errorf("create resource: %w", err)
		}

		// If the resource is not yet ready, start async polling
		if providerOutput.Status != "active" && providerOutput.Status != "running" && providerOutput.Status != "Ready" {
			providerResourceID := providerOutput.ID
			s.startAsyncStatusPoller(id, providerResourceID, p)
		} else {
			if err := s.resourceRepo.UpdateResource(ctx, id, values, "active"); err != nil {
				return nil, fmt.Errorf("activate resource: %w", err)
			}
			res.Status = "active"
		}

		return res, nil
	}

	// Showbiz-managed resource: no connection needed
	id = fmt.Sprintf("sbz:%s:%s:%s", input.ResourceType, projectID, input.Name)

	values := input.Values
	if values == nil {
		values = make(map[string]interface{})
	}

	// For networks, compute gateway from CIDR
	if input.ResourceType == "network" {
		if cidrStr, ok := values["cidr"].(string); ok {
			ip, ipNet, err := net.ParseCIDR(cidrStr)
			if err == nil {
				gw := make(net.IP, len(ip))
				copy(gw, ipNet.IP)
				// Gateway is the first usable IP (network address + 1)
				gw[len(gw)-1]++
				values["gateway"] = gw.String()
			}
		}
	}

	res := &model.Resource{
		ID:           id,
		Name:         input.Name,
		ProjectID:    projectID,
		ConnectionID: nil,
		ResourceType: input.ResourceType,
		Values:       values,
		Status:       "active",
	}

	if err := s.resourceRepo.CreateResource(ctx, res); err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	return res, nil
}

// startAsyncStatusPoller polls the provider every second until the resource reaches Ready status.
func (s *ResourceService) startAsyncStatusPoller(resourceID, providerResourceID string, p provider.Provider) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		timeout := time.After(5 * time.Minute)

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				output, err := p.GetResource(ctx, providerResourceID)
				cancel()

				if err != nil {
					slog.Warn("polling provider resource failed", "resourceID", resourceID, "error", err)
					continue
				}

				slog.Info("polling resource status", "resourceID", resourceID, "status", output.Status)

				if output.Status == "Ready" || output.Status == "running" {
					// Resource is ready — update DB with final values and status
					values := output.Properties
					if values == nil {
						values = make(map[string]interface{})
					}
					ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
					if err := s.resourceRepo.UpdateResource(ctx2, resourceID, values, "active"); err != nil {
						slog.Error("failed to update resource to active", "resourceID", resourceID, "error", err)
					}
					cancel2()
					return
				}

				if output.Status == "Failed" {
					ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
					if err := s.resourceRepo.UpdateResource(ctx2, resourceID, output.Properties, "failed"); err != nil {
						slog.Error("failed to update resource to failed", "resourceID", resourceID, "error", err)
					}
					cancel2()
					return
				}

			case <-timeout:
				slog.Error("resource polling timed out", "resourceID", resourceID)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_ = s.resourceRepo.UpdateResource(ctx, resourceID, nil, "failed")
				cancel()
				return
			}
		}
	}()
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

	// Call provider to delete the resource if it has a connection
	if res.ConnectionID != nil {
		conn, err := s.connRepo.GetConnectionByID(ctx, *res.ConnectionID)
		if err == nil && conn != nil {
			if p, ok := s.registry.Get(conn.Provider); ok {
				if providerID, ok := res.Values["fakeproviderID"]; ok {
					if idStr, ok := providerID.(string); ok {
						_ = p.DeleteResource(ctx, idStr)
					}
				}
			}
		}
	}

	if err := s.resourceRepo.DeleteResource(ctx, id); err != nil {
		return fmt.Errorf("delete resource: %w", err)
	}
	return nil
}
