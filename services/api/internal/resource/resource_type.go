package resource

// ResourceType defines the contract every resource type must implement.
type ResourceType interface {
	// Name returns the type identifier (e.g., "machine", "network").
	Name() string

	// RequiresConnection returns true if this type needs a provider connection.
	RequiresConnection() bool

	// ValidateCreate validates input values before creating a resource.
	ValidateCreate(values map[string]interface{}) error

	// ValidateUpdate validates values before updating an existing resource.
	ValidateUpdate(currentValues, newValues map[string]interface{}) error

	// InputSchema returns the schema of fields accepted during creation.
	InputSchema() []FieldSchema

	// OutputSchema returns the schema of fields returned after provisioning.
	OutputSchema() []FieldSchema
}

// FieldSchema describes a single field in a resource type's input or output.
type FieldSchema struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "string", "number", "boolean"
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// ResourceTypeInfo contains metadata about a registered resource type, used for API responses.
type ResourceTypeInfo struct {
	Name               string        `json:"name"`
	RequiresConnection bool          `json:"requiresConnection"`
	InputSchema        []FieldSchema `json:"inputSchema"`
	OutputSchema       []FieldSchema `json:"outputSchema"`
}
