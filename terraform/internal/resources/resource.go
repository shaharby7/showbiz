package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"

	"github.com/shaharby7/showbiz/terraform/internal/helpers"
)

var (
	_ resource.Resource                = &resourceResource{}
	_ resource.ResourceWithConfigure   = &resourceResource{}
	_ resource.ResourceWithImportState = &resourceResource{}
)

type resourceResource struct {
	client *showbiz.Client
}

type resourceResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectID    types.String `tfsdk:"project_id"`
	ConnectionID types.String `tfsdk:"connection_id"`
	Name         types.String `tfsdk:"name"`
	ResourceType types.String `tfsdk:"resource_type"`
	Values       types.Map    `tfsdk:"values"`
	Status       types.String `tfsdk:"status"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func NewResourceResource() resource.Resource {
	return &resourceResource{}
}

func (r *resourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (r *resourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Showbiz resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID this resource belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connection_id": schema.StringAttribute{
				Description: "The connection ID used to manage this resource.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Description: "The type of resource (e.g. machine, network).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"values": schema.MapAttribute{
				Description: "Resource configuration values.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"status": schema.StringAttribute{
				Description: "The current status of the resource.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The last update timestamp.",
				Computed:    true,
			},
		},
	}
}

func (r *resourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*showbiz.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *showbiz.Client, got unexpected type.",
		)
		return
	}
	r.client = client
}

func (r *resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	values, diags := helpers.ExpandStringMap(ctx, plan.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.CreateResourceInput{
		Name:         plan.Name.ValueString(),
		ConnectionID: plan.ConnectionID.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		Values:       values,
	}

	res, err := r.client.CreateResource(ctx, plan.ProjectID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}

	plan.ID = types.StringValue(res.ID)
	plan.ProjectID = types.StringValue(res.ProjectID)
	plan.ConnectionID = types.StringValue(res.ConnectionID)
	plan.Name = types.StringValue(res.Name)
	plan.ResourceType = types.StringValue(res.ResourceType)
	plan.Status = types.StringValue(res.Status)

	valuesMap, diags := helpers.FlattenStringMap(res.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Values = valuesMap

	plan.CreatedAt = types.StringValue(res.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(res.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetResource(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if showbiz.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	state.ID = types.StringValue(res.ID)
	state.ProjectID = types.StringValue(res.ProjectID)
	state.ConnectionID = types.StringValue(res.ConnectionID)
	state.Name = types.StringValue(res.Name)
	state.ResourceType = types.StringValue(res.ResourceType)
	state.Status = types.StringValue(res.Status)

	valuesMap, diags := helpers.FlattenStringMap(res.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Values = valuesMap

	state.CreatedAt = types.StringValue(res.CreatedAt.String())
	state.UpdatedAt = types.StringValue(res.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	values, diags := helpers.ExpandStringMap(ctx, plan.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.UpdateResourceInput{
		Values: values,
	}

	res, err := r.client.UpdateResource(ctx, state.ProjectID.ValueString(), state.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	plan.ID = types.StringValue(res.ID)
	plan.ProjectID = types.StringValue(res.ProjectID)
	plan.ConnectionID = types.StringValue(res.ConnectionID)
	plan.Name = types.StringValue(res.Name)
	plan.ResourceType = types.StringValue(res.ResourceType)
	plan.Status = types.StringValue(res.Status)

	valuesMap, diags := helpers.FlattenStringMap(res.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Values = valuesMap

	plan.CreatedAt = types.StringValue(res.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(res.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteResource(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}
}

func (r *resourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split by first "/" only — resource IDs may contain colons.
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: <project_id>/<resource_id>, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
