package resource

import "sync"

// Registry holds all registered resource types.
type Registry struct {
	mu    sync.RWMutex
	types map[string]ResourceType
}

// NewRegistry creates a new resource type registry.
func NewRegistry() *Registry {
	return &Registry{
		types: make(map[string]ResourceType),
	}
}

// Register adds a resource type to the registry.
func (r *Registry) Register(rt ResourceType) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.types[rt.Name()] = rt
}

// Get returns a resource type by name.
func (r *Registry) Get(name string) (ResourceType, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, ok := r.types[name]
	return rt, ok
}

// List returns metadata for all registered resource types.
func (r *Registry) List() []ResourceTypeInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]ResourceTypeInfo, 0, len(r.types))
	for _, rt := range r.types {
		infos = append(infos, ResourceTypeInfo{
			Name:               rt.Name(),
			RequiresConnection: rt.RequiresConnection(),
			InputSchema:        rt.InputSchema(),
			OutputSchema:       rt.OutputSchema(),
		})
	}
	return infos
}
