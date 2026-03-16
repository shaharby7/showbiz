package resource

import "fmt"

// MachineResourceType implements ResourceType for compute instances.
type MachineResourceType struct{}

var _ ResourceType = (*MachineResourceType)(nil)

func NewMachineResourceType() *MachineResourceType {
	return &MachineResourceType{}
}

func (m *MachineResourceType) Name() string { return "machine" }

func (m *MachineResourceType) RequiresConnection() bool { return true }

func (m *MachineResourceType) ValidateCreate(values map[string]interface{}) error {
	if values == nil {
		return fmt.Errorf("values are required")
	}
	cpu, ok := values["cpu"]
	if !ok {
		return fmt.Errorf("cpu is required")
	}
	if n, ok := cpu.(float64); !ok || n < 1 {
		return fmt.Errorf("cpu must be a positive number")
	}

	mem, ok := values["memoryMB"]
	if !ok {
		return fmt.Errorf("memoryMB is required")
	}
	if n, ok := mem.(float64); !ok || n < 1 {
		return fmt.Errorf("memoryMB must be a positive number")
	}

	image, ok := values["image"]
	if !ok {
		return fmt.Errorf("image is required")
	}
	if s, ok := image.(string); !ok || s == "" {
		return fmt.Errorf("image must be a non-empty string")
	}

	if ns, ok := values["namespace"]; ok {
		if _, ok := ns.(string); !ok {
			return fmt.Errorf("namespace must be a string")
		}
	}

	return nil
}

func (m *MachineResourceType) ValidateUpdate(_, newValues map[string]interface{}) error {
	return m.ValidateCreate(newValues)
}

func (m *MachineResourceType) InputSchema() []FieldSchema {
	return []FieldSchema{
		{Name: "cpu", Type: "number", Required: true, Description: "Number of CPU cores"},
		{Name: "memoryMB", Type: "number", Required: true, Description: "Memory in megabytes"},
		{Name: "image", Type: "string", Required: true, Description: "OS image identifier"},
		{Name: "namespace", Type: "string", Required: false, Description: "Target namespace"},
	}
}

func (m *MachineResourceType) OutputSchema() []FieldSchema {
	return []FieldSchema{
		{Name: "ip", Type: "string", Description: "Assigned IP address"},
		{Name: "providerResourceId", Type: "string", Description: "Provider-side resource ID"},
	}
}
