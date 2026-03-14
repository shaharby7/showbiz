package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
	"github.com/showbiz-io/showbiz/services/api/internal/repository"
)

type OrgService struct {
	orgRepo    *repository.OrgRepo
	memberRepo *repository.MemberRepo
}

func NewOrgService(orgRepo *repository.OrgRepo, memberRepo *repository.MemberRepo) *OrgService {
	return &OrgService{
		orgRepo:    orgRepo,
		memberRepo: memberRepo,
	}
}

type CreateOrgInput struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type UpdateOrgInput struct {
	DisplayName string `json:"displayName"`
}

func (s *OrgService) Create(ctx context.Context, input CreateOrgInput) (*model.Organization, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.DisplayName == "" {
		input.DisplayName = input.Name
	}

	org := &model.Organization{
		ID:          uuid.New().String(),
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Active:      true,
	}

	if err := s.orgRepo.CreateOrganization(ctx, org); err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	created, err := s.orgRepo.GetOrganizationByID(ctx, org.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created organization: %w", err)
	}
	return created, nil
}

func (s *OrgService) Get(ctx context.Context, id string) (*model.Organization, error) {
	org, err := s.orgRepo.GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}
	return org, nil
}

func (s *OrgService) Update(ctx context.Context, id string, input UpdateOrgInput) (*model.Organization, error) {
	if input.DisplayName == "" {
		return nil, fmt.Errorf("displayName is required")
	}

	if err := s.orgRepo.UpdateOrganization(ctx, id, input.DisplayName); err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}

	org, err := s.orgRepo.GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch updated organization: %w", err)
	}
	return org, nil
}

func (s *OrgService) Deactivate(ctx context.Context, id string) error {
	org, err := s.orgRepo.GetOrganizationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get organization: %w", err)
	}
	if org == nil {
		return fmt.Errorf("organization not found")
	}

	if err := s.orgRepo.DeleteProjectsByOrgID(ctx, id); err != nil {
		return fmt.Errorf("delete projects: %w", err)
	}

	if err := s.orgRepo.SetOrganizationActive(ctx, id, false); err != nil {
		return fmt.Errorf("deactivate organization: %w", err)
	}
	return nil
}

func (s *OrgService) Activate(ctx context.Context, id string) error {
	org, err := s.orgRepo.GetOrganizationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get organization: %w", err)
	}
	if org == nil {
		return fmt.Errorf("organization not found")
	}

	if err := s.orgRepo.SetOrganizationActive(ctx, id, true); err != nil {
		return fmt.Errorf("activate organization: %w", err)
	}
	return nil
}

func (s *OrgService) List(ctx context.Context, cursor string, limit int) ([]*model.Organization, string, error) {
	orgs, nextCursor, err := s.orgRepo.ListOrganizations(ctx, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("list organizations: %w", err)
	}
	return orgs, nextCursor, nil
}

func (s *OrgService) ListMembers(ctx context.Context, orgID string) ([]*model.User, error) {
	members, err := s.memberRepo.ListMembers(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return members, nil
}

func (s *OrgService) AddMember(ctx context.Context, orgID, email string) error {
	if err := s.memberRepo.AddMember(ctx, orgID, email); err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (s *OrgService) RemoveMember(ctx context.Context, orgID, email string) error {
	if err := s.memberRepo.RemoveMember(ctx, orgID, email); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}
