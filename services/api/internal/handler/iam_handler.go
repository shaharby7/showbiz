package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/showbiz-io/showbiz/services/api/internal/service"
)

type IAMHandler struct {
	svc *service.IAMService
}

func NewIAMHandler(svc *service.IAMService) *IAMHandler {
	return &IAMHandler{svc: svc}
}

func (h *IAMHandler) ListGlobalPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.svc.ListGlobalPolicies(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list global policies")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": policies,
	})
}

func (h *IAMHandler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "policyId")

	policy, err := h.svc.GetPolicy(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "policy not found") {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Policy not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get policy")
		return
	}

	JSON(w, http.StatusOK, policy)
}

func (h *IAMHandler) ListOrgPolicies(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	policies, err := h.svc.ListOrgPolicies(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list organization policies")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": policies,
	})
}

func (h *IAMHandler) CreateOrgPolicy(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	var input service.CreatePolicyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Name == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}
	if len(input.Permissions) == 0 {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "At least one permission is required")
		return
	}

	policy, err := h.svc.CreateOrgPolicy(r.Context(), orgID, input)
	if err != nil {
		if strings.Contains(err.Error(), "is required") || strings.Contains(err.Error(), "invalid permission") {
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create policy")
		return
	}

	JSON(w, http.StatusCreated, policy)
}

func (h *IAMHandler) DeleteOrgPolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "policyId")

	err := h.svc.DeleteOrgPolicy(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "policy not found") {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Policy not found")
			return
		}
		if strings.Contains(err.Error(), "cannot delete global policy") {
			Error(w, http.StatusForbidden, "FORBIDDEN", "Cannot delete global policy")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete policy")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}

func (h *IAMHandler) ListProjectAttachments(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	attachments, err := h.svc.ListProjectAttachments(r.Context(), projectID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list attachments")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": attachments,
	})
}

func (h *IAMHandler) AttachPolicy(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	var body struct {
		UserEmail string `json:"userEmail"`
		PolicyID  string `json:"policyId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if body.UserEmail == "" || body.PolicyID == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "userEmail and policyId are required")
		return
	}

	input := service.AttachPolicyInput{
		ProjectID: projectID,
		UserEmail: body.UserEmail,
		PolicyID:  body.PolicyID,
	}

	attachment, err := h.svc.AttachPolicy(r.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "are required") {
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to attach policy")
		return
	}

	JSON(w, http.StatusCreated, attachment)
}

func (h *IAMHandler) DetachPolicy(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	var body struct {
		UserEmail string `json:"userEmail"`
		PolicyID  string `json:"policyId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if body.UserEmail == "" || body.PolicyID == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "userEmail and policyId are required")
		return
	}

	err := h.svc.DetachPolicy(r.Context(), projectID, body.UserEmail, body.PolicyID)
	if err != nil {
		if strings.Contains(err.Error(), "attachment not found") {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Attachment not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to detach policy")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}
