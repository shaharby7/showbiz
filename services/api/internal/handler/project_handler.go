package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/showbiz-io/showbiz/services/api/internal/service"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	var input service.CreateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Name == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}

	project, err := h.svc.Create(r.Context(), orgID, input)
	if err != nil {
		switch err.Error() {
		case "name is required":
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		case "organization not found":
			Error(w, http.StatusNotFound, "NOT_FOUND", "Organization not found")
		case "organization is not active":
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Organization is not active")
		case "project name already exists":
			Error(w, http.StatusConflict, "CONFLICT", "Project name already exists in this organization")
		default:
			Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create project")
		}
		return
	}

	JSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	project, err := h.svc.Get(r.Context(), projectID)
	if err != nil {
		if err.Error() == "project not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get project")
		return
	}

	JSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	projects, nextCursor, err := h.svc.List(r.Context(), orgID, cursor, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list projects")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": projects,
		"pagination": map[string]interface{}{
			"nextCursor": nextCursor,
			"hasMore":    nextCursor != "",
		},
	})
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	var input service.UpdateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	project, err := h.svc.Update(r.Context(), projectID, input)
	if err != nil {
		if err.Error() == "project not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update project")
		return
	}

	JSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	err := h.svc.Delete(r.Context(), projectID)
	if err != nil {
		if err.Error() == "project not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete project")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}
