package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/shaharby7/showbiz/services/api/internal/service"
)

type ResourceHandler struct {
	svc *service.ResourceService
}

func NewResourceHandler(svc *service.ResourceService) *ResourceHandler {
	return &ResourceHandler{svc: svc}
}

func (h *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	var input service.CreateResourceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Name == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}
	if input.ResourceType == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "ResourceType is required")
		return
	}

	res, err := h.svc.Create(r.Context(), projectID, input)
	if err != nil {
		switch err.Error() {
		case "connection not found":
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Connection not found")
		case "provider not found":
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Provider not found")
		case "invalid resource type":
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid resource type for this connection's provider")
		case "unknown resource type":
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Unknown resource type")
		case "resource name already exists":
			Error(w, http.StatusConflict, "CONFLICT", "Resource name already exists in this project")
		default:
			if len(err.Error()) > 17 && err.Error()[:17] == "validation error:" {
				Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			} else if len(err.Error()) > 24 && err.Error()[:24] == "connectionId is required" {
				Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			} else {
				Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create resource")
			}
		}
		return
	}

	JSON(w, http.StatusCreated, res)
}

func (h *ResourceHandler) Get(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "resourceId")
	id, err := url.PathUnescape(rawID)
	if err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid resource ID")
		return
	}

	res, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err.Error() == "resource not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get resource")
		return
	}

	JSON(w, http.StatusOK, res)
}

func (h *ResourceHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	resources, nextCursor, err := h.svc.List(r.Context(), projectID, cursor, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list resources")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": resources,
		"pagination": map[string]interface{}{
			"nextCursor": nextCursor,
			"hasMore":    nextCursor != "",
		},
	})
}

func (h *ResourceHandler) Update(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "resourceId")
	id, err := url.PathUnescape(rawID)
	if err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid resource ID")
		return
	}

	var input service.UpdateResourceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Values == nil {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Values is required")
		return
	}

	res, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		if err.Error() == "resource not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update resource")
		return
	}

	JSON(w, http.StatusOK, res)
}

func (h *ResourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "resourceId")
	id, err := url.PathUnescape(rawID)
	if err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid resource ID")
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "resource not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete resource")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}
