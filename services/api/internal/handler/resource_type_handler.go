package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shaharby7/showbiz/services/api/internal/resource"
)

// ResourceTypeHandler handles resource-type-related HTTP requests.
type ResourceTypeHandler struct {
	registry *resource.Registry
}

// NewResourceTypeHandler creates a new ResourceTypeHandler.
func NewResourceTypeHandler(registry *resource.Registry) *ResourceTypeHandler {
	return &ResourceTypeHandler{registry: registry}
}

// List returns all registered resource types with their schemas.
func (h *ResourceTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, h.registry.List())
}

// Get returns a single resource type by name.
func (h *ResourceTypeHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rt, ok := h.registry.Get(name)
	if !ok {
		Error(w, http.StatusNotFound, "RESOURCE_TYPE_NOT_FOUND", "Resource type not found")
		return
	}

	JSON(w, http.StatusOK, resource.ResourceTypeInfo{
		Name:               rt.Name(),
		RequiresConnection: rt.RequiresConnection(),
		InputSchema:        rt.InputSchema(),
		OutputSchema:       rt.OutputSchema(),
	})
}
