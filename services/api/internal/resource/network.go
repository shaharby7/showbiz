package resource

import (
	"fmt"
	"net"
)

// NetworkResourceType implements ResourceType for Showbiz-managed virtual networks.
type NetworkResourceType struct{}

var _ ResourceType = (*NetworkResourceType)(nil)

func NewNetworkResourceType() *NetworkResourceType {
	return &NetworkResourceType{}
}

func (n *NetworkResourceType) Name() string { return "network" }

func (n *NetworkResourceType) RequiresConnection() bool { return false }

func (n *NetworkResourceType) ValidateCreate(values map[string]interface{}) error {
	if values == nil {
		return fmt.Errorf("values are required")
	}
	cidr, ok := values["cidr"]
	if !ok {
		return fmt.Errorf("cidr is required")
	}
	cidrStr, ok := cidr.(string)
	if !ok || cidrStr == "" {
		return fmt.Errorf("cidr must be a non-empty string")
	}
	_, _, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return fmt.Errorf("cidr must be valid CIDR notation (e.g., 10.0.0.0/16)")
	}

	if desc, ok := values["description"]; ok {
		if _, ok := desc.(string); !ok {
			return fmt.Errorf("description must be a string")
		}
	}

	return nil
}

func (n *NetworkResourceType) ValidateUpdate(_, newValues map[string]interface{}) error {
	if newValues == nil {
		return fmt.Errorf("values are required")
	}
	if desc, ok := newValues["description"]; ok {
		if _, ok := desc.(string); !ok {
			return fmt.Errorf("description must be a string")
		}
	}
	return nil
}

func (n *NetworkResourceType) InputSchema() []FieldSchema {
	return []FieldSchema{
		{Name: "cidr", Type: "string", Required: true, Description: "Network CIDR block (e.g., 10.0.0.0/16)"},
		{Name: "description", Type: "string", Required: false, Description: "Network description"},
	}
}

func (n *NetworkResourceType) OutputSchema() []FieldSchema {
	return []FieldSchema{
		{Name: "gateway", Type: "string", Description: "Network gateway address"},
	}
}
