package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExpandStringMap converts a types.Map to map[string]interface{}.
func ExpandStringMap(ctx context.Context, m types.Map) (map[string]interface{}, diag.Diagnostics) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}
	elements := make(map[string]types.String, len(m.Elements()))
	diags := m.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make(map[string]interface{}, len(elements))
	for k, v := range elements {
		result[k] = v.ValueString()
	}
	return result, nil
}

// FlattenStringMap converts map[string]interface{} to a types.Map.
func FlattenStringMap(m map[string]interface{}) (types.Map, diag.Diagnostics) {
	if m == nil {
		return types.MapNull(types.StringType), nil
	}
	elements := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elements[k] = types.StringValue(fmt.Sprintf("%v", v))
	}
	return types.MapValue(types.StringType, elements)
}
