package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
	"github.com/showbiz-io/showbiz/services/api/internal/repository"
)

type IAMService struct {
	policyRepo     *repository.PolicyRepo
	attachmentRepo *repository.PolicyAttachmentRepo
	projectRepo    *repository.ProjectRepo
	userRepo       *repository.UserRepo
}

func NewIAMService(
	policyRepo *repository.PolicyRepo,
	attachmentRepo *repository.PolicyAttachmentRepo,
	projectRepo *repository.ProjectRepo,
	userRepo *repository.UserRepo,
) *IAMService {
	return &IAMService{
		policyRepo:     policyRepo,
		attachmentRepo: attachmentRepo,
		projectRepo:    projectRepo,
		userRepo:       userRepo,
	}
}

type CreatePolicyInput struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type AttachPolicyInput struct {
	ProjectID string `json:"projectId"`
	UserEmail string `json:"userEmail"`
	PolicyID  string `json:"policyId"`
}

func (s *IAMService) CreateOrgPolicy(ctx context.Context, orgID string, input CreatePolicyInput) (*model.Policy, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if len(input.Permissions) == 0 {
		return nil, fmt.Errorf("at least one permission is required")
	}
	for _, perm := range input.Permissions {
		if !isValidPermission(perm) {
			return nil, fmt.Errorf("invalid permission format: %s", perm)
		}
	}

	policy := &model.Policy{
		ID:             uuid.New().String(),
		Name:           input.Name,
		Scope:          "organization",
		OrganizationID: orgID,
		Permissions:    input.Permissions,
	}

	if err := s.policyRepo.CreatePolicy(ctx, policy); err != nil {
		return nil, fmt.Errorf("create policy: %w", err)
	}

	created, err := s.policyRepo.GetPolicyByID(ctx, policy.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created policy: %w", err)
	}
	return created, nil
}

func (s *IAMService) GetPolicy(ctx context.Context, id string) (*model.Policy, error) {
	policy, err := s.policyRepo.GetPolicyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get policy: %w", err)
	}
	if policy == nil {
		return nil, fmt.Errorf("policy not found")
	}
	return policy, nil
}

func (s *IAMService) ListGlobalPolicies(ctx context.Context) ([]*model.Policy, error) {
	policies, err := s.policyRepo.ListGlobalPolicies(ctx)
	if err != nil {
		return nil, fmt.Errorf("list global policies: %w", err)
	}
	return policies, nil
}

func (s *IAMService) ListOrgPolicies(ctx context.Context, orgID string) ([]*model.Policy, error) {
	policies, err := s.policyRepo.ListOrgPolicies(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("list org policies: %w", err)
	}
	return policies, nil
}

func (s *IAMService) DeleteOrgPolicy(ctx context.Context, id string) error {
	policy, err := s.policyRepo.GetPolicyByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get policy: %w", err)
	}
	if policy == nil {
		return fmt.Errorf("policy not found")
	}
	if policy.Scope == "global" {
		return fmt.Errorf("cannot delete global policy")
	}

	if err := s.policyRepo.DeletePolicy(ctx, id); err != nil {
		return fmt.Errorf("delete policy: %w", err)
	}
	return nil
}

func (s *IAMService) AttachPolicy(ctx context.Context, input AttachPolicyInput) (*model.PolicyAttachment, error) {
	if input.ProjectID == "" || input.UserEmail == "" || input.PolicyID == "" {
		return nil, fmt.Errorf("projectId, userEmail, and policyId are required")
	}

	policy, err := s.policyRepo.GetPolicyByID(ctx, input.PolicyID)
	if err != nil {
		return nil, fmt.Errorf("get policy: %w", err)
	}
	if policy == nil {
		return nil, fmt.Errorf("policy not found")
	}

	project, err := s.projectRepo.GetProjectByID(ctx, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, input.UserEmail)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	attachment := &model.PolicyAttachment{
		ID:        uuid.New().String(),
		ProjectID: input.ProjectID,
		UserEmail: input.UserEmail,
		PolicyID:  input.PolicyID,
	}

	if err := s.attachmentRepo.AttachPolicy(ctx, attachment); err != nil {
		return nil, fmt.Errorf("attach policy: %w", err)
	}
	return attachment, nil
}

func (s *IAMService) DetachPolicy(ctx context.Context, projectID, userEmail, policyID string) error {
	if err := s.attachmentRepo.DetachPolicy(ctx, projectID, userEmail, policyID); err != nil {
		return fmt.Errorf("detach policy: %w", err)
	}
	return nil
}

func (s *IAMService) ListProjectAttachments(ctx context.Context, projectID string) ([]*model.PolicyAttachment, error) {
	attachments, err := s.attachmentRepo.ListAttachmentsByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list project attachments: %w", err)
	}
	return attachments, nil
}

func (s *IAMService) CheckPermission(ctx context.Context, userEmail, projectID, entity, action string) (bool, error) {
	perms, err := s.attachmentRepo.GetUserPermissions(ctx, userEmail, projectID)
	if err != nil {
		return false, fmt.Errorf("get user permissions: %w", err)
	}

	for _, perm := range perms {
		if perm == "*:*" {
			return true, nil
		}
		if perm == entity+":"+action {
			return true, nil
		}
		parts := strings.SplitN(perm, ":", 2)
		if len(parts) == 2 {
			if parts[0] == "*" && parts[1] == action {
				return true, nil
			}
			if parts[0] == entity && parts[1] == "*" {
				return true, nil
			}
		}
	}
	return false, nil
}

func isValidPermission(perm string) bool {
	parts := strings.SplitN(perm, ":", 2)
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}
