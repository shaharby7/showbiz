package provider

import "sync"

// ProviderInfo contains metadata about a registered provider.
type ProviderInfo struct {
	Name          string   `json:"name"`
	ResourceTypes []string `json:"resourceTypes"`
}

// Registry holds all registered providers.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// NewRegistry creates a new provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(name string, p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = p
}

// Get returns a provider by name.
func (r *Registry) Get(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	return p, ok
}

// List returns metadata for all registered providers.
func (r *Registry) List() []ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]ProviderInfo, 0, len(r.providers))
	for _, p := range r.providers {
		infos = append(infos, ProviderInfo{
			Name:          p.Name(),
			ResourceTypes: p.ResourceTypes(),
		})
	}
	return infos
}
