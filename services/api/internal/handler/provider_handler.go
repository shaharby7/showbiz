package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/showbiz-io/showbiz/services/api/internal/provider"
)

// ProviderHandler handles provider-related HTTP requests.
type ProviderHandler struct {
	registry *provider.Registry
}

// NewProviderHandler creates a new ProviderHandler.
func NewProviderHandler(registry *provider.Registry) *ProviderHandler {
	return &ProviderHandler{registry: registry}
}

// List returns all registered providers.
func (h *ProviderHandler) List(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, h.registry.List())
}

// Get returns a single provider by name.
func (h *ProviderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, ok := h.registry.Get(id)
	if !ok {
		Error(w, http.StatusNotFound, "PROVIDER_NOT_FOUND", "Provider not found")
		return
	}

	JSON(w, http.StatusOK, provider.ProviderInfo{
		Name:          p.Name(),
		ResourceTypes: p.ResourceTypes(),
	})
}
