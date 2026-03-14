package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/showbiz-io/showbiz/services/api/internal/service"
)

type OrgHandler struct {
	svc *service.OrgService
}

func NewOrgHandler(svc *service.OrgService) *OrgHandler {
	return &OrgHandler{svc: svc}
}

func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input service.CreateOrgInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Name == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}

	org, err := h.svc.Create(r.Context(), input)
	if err != nil {
		if err.Error() == "name is required" {
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create organization")
		return
	}

	JSON(w, http.StatusCreated, org)
}

func (h *OrgHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	org, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err.Error() == "organization not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Organization not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get organization")
		return
	}

	JSON(w, http.StatusOK, org)
}

func (h *OrgHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input service.UpdateOrgInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.DisplayName == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "DisplayName is required")
		return
	}

	org, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		if err.Error() == "organization not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Organization not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update organization")
		return
	}

	JSON(w, http.StatusOK, org)
}

func (h *OrgHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.svc.Deactivate(r.Context(), id)
	if err != nil {
		if err.Error() == "organization not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Organization not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to deactivate organization")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}

func (h *OrgHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.svc.Activate(r.Context(), id)
	if err != nil {
		if err.Error() == "organization not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Organization not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to activate organization")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}

func (h *OrgHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	orgs, nextCursor, err := h.svc.List(r.Context(), cursor, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list organizations")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": orgs,
		"pagination": map[string]interface{}{
			"nextCursor": nextCursor,
			"hasMore":    nextCursor != "",
		},
	})
}

func (h *OrgHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")

	members, err := h.svc.ListMembers(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list members")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": members,
	})
}

func (h *OrgHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")

	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Email == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email is required")
		return
	}

	err := h.svc.AddMember(r.Context(), orgID, input.Email)
	if err != nil {
		if err.Error() == "add member: user not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add member")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}

func (h *OrgHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	email := chi.URLParam(r, "email")

	err := h.svc.RemoveMember(r.Context(), orgID, email)
	if err != nil {
		if err.Error() == "remove member: member not found in organization" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Member not found in organization")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove member")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}
