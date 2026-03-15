package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/shaharby7/showbiz/services/api/internal/service"
)

type ConnectionHandler struct {
	svc *service.ConnectionService
}

func NewConnectionHandler(svc *service.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{svc: svc}
}

func (h *ConnectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")

	var input service.CreateConnectionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Name == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}
	if input.Provider == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Provider is required")
		return
	}

	conn, err := h.svc.Create(r.Context(), projectID, input)
	if err != nil {
		if err.Error() == "provider not found" {
			Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Provider not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create connection")
		return
	}

	JSON(w, http.StatusCreated, conn)
}

func (h *ConnectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "connectionId")

	conn, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err.Error() == "connection not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Connection not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get connection")
		return
	}

	JSON(w, http.StatusOK, conn)
}

func (h *ConnectionHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	conns, nextCursor, err := h.svc.List(r.Context(), projectID, cursor, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list connections")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": conns,
		"pagination": map[string]interface{}{
			"nextCursor": nextCursor,
			"hasMore":    nextCursor != "",
		},
	})
}

func (h *ConnectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "connectionId")

	var input service.UpdateConnectionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Config == nil {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Config is required")
		return
	}

	conn, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		if err.Error() == "connection not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Connection not found")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update connection")
		return
	}

	JSON(w, http.StatusOK, conn)
}

func (h *ConnectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "connectionId")

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "connection not found" {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Connection not found")
			return
		}
		if err.Error() == "connection has resources" {
			Error(w, http.StatusConflict, "CONFLICT", "Connection has resources and cannot be deleted")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete connection")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}
