package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaharby7/showbiz/services/api/internal/model"
	"github.com/shaharby7/showbiz/services/api/internal/repository"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepo
	orgRepo     *repository.OrgRepo
}

func NewProjectService(projectRepo *repository.ProjectRepo, orgRepo *repository.OrgRepo) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		orgRepo:     orgRepo,
	}
}

type CreateProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateProjectInput struct {
	Description string `json:"description"`
}

func (s *ProjectService) Create(ctx context.Context, orgID string, input CreateProjectInput) (*model.Project, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	org, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}
	if !org.Active {
		return nil, fmt.Errorf("organization is not active")
	}

	existing, err := s.projectRepo.GetProjectByOrgIDAndName(ctx, orgID, input.Name)
	if err != nil {
		return nil, fmt.Errorf("check project name: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("project name already exists")
	}

	project := &model.Project{
		ID:             uuid.New().String(),
		Name:           input.Name,
		OrganizationID: orgID,
		Description:    input.Description,
	}

	if err := s.projectRepo.CreateProject(ctx, project); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	created, err := s.projectRepo.GetProjectByID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created project: %w", err)
	}
	return created, nil
}

func (s *ProjectService) Get(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}
	return project, nil
}

func (s *ProjectService) List(ctx context.Context, orgID string, cursor string, limit int) ([]*model.Project, string, error) {
	projects, nextCursor, err := s.projectRepo.ListProjectsByOrgID(ctx, orgID, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("list projects: %w", err)
	}
	return projects, nextCursor, nil
}

func (s *ProjectService) Update(ctx context.Context, id string, input UpdateProjectInput) (*model.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	if err := s.projectRepo.UpdateProject(ctx, id, input.Description); err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	updated, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch updated project: %w", err)
	}
	return updated, nil
}

func (s *ProjectService) Delete(ctx context.Context, id string) error {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	if err := s.projectRepo.DeleteProject(ctx, id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
