package showbiz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListGlobalPolicies retrieves all global IAM policies.
func (c *Client) ListGlobalPolicies(ctx context.Context) ([]*Policy, error) {
	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, "/v1/iam/policies", nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data []*Policy `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode global policies: %w", err)
	}

	return envelope.Data, nil
}

// GetPolicy retrieves a policy by ID.
func (c *Client) GetPolicy(ctx context.Context, policyID string) (*Policy, error) {
	var policy Policy
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/iam/policies/%s", policyID), nil, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

// ListOrgPolicies retrieves policies for an organization.
func (c *Client) ListOrgPolicies(ctx context.Context, orgID string) ([]*Policy, error) {
	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s/policies", orgID), nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data []*Policy `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode org policies: %w", err)
	}

	return envelope.Data, nil
}

// CreateOrgPolicy creates a new policy in an organization.
func (c *Client) CreateOrgPolicy(ctx context.Context, orgID string, input CreatePolicyInput) (*Policy, error) {
	var policy Policy
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/policies", orgID), input, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

// DeleteOrgPolicy deletes a policy from an organization.
func (c *Client) DeleteOrgPolicy(ctx context.Context, orgID, policyID string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/policies/%s", orgID, policyID), nil, nil)
}

// ListPolicyAttachments retrieves policy attachments for a project.
func (c *Client) ListPolicyAttachments(ctx context.Context, orgID, projectID string) ([]*PolicyAttachment, error) {
	var raw json.RawMessage
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s/projects/%s/attachments", orgID, projectID), nil, &raw); err != nil {
		return nil, err
	}

	var envelope struct {
		Data []*PolicyAttachment `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("showbiz: failed to decode policy attachments: %w", err)
	}

	return envelope.Data, nil
}

// AttachPolicy attaches a policy to a user in a project.
func (c *Client) AttachPolicy(ctx context.Context, orgID, projectID string, input AttachPolicyInput) (*PolicyAttachment, error) {
	var attachment PolicyAttachment
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/projects/%s/attachments", orgID, projectID), input, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}

// DetachPolicy detaches a policy from a user in a project.
func (c *Client) DetachPolicy(ctx context.Context, orgID, projectID string, input DetachPolicyInput) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/projects/%s/attachments", orgID, projectID), input, nil)
}
