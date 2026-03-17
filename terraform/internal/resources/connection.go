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
	_ resource.Resource                = &connectionResource{}
	_ resource.ResourceWithConfigure   = &connectionResource{}
	_ resource.ResourceWithImportState = &connectionResource{}
)

type connectionResource struct {
	client *showbiz.Client
}

type connectionResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Name         types.String `tfsdk:"name"`
	ProviderName types.String `tfsdk:"provider_name"`
	Credentials  types.Map    `tfsdk:"credentials"`
	Config       types.Map    `tfsdk:"config"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func NewConnectionResource() resource.Resource {
	return &connectionResource{}
}

func (r *connectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (r *connectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Showbiz connection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the connection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID this connection belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_name": schema.StringAttribute{
				Description: "The cloud provider type (e.g. aws, gcp).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"credentials": schema.MapAttribute{
				Description: "Provider credentials (write-only, not returned by API).",
				Optional:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"config": schema.MapAttribute{
				Description: "Provider configuration.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
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

func (r *connectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *connectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan connectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds, diags := helpers.ExpandStringMap(ctx, plan.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, diags := helpers.ExpandStringMap(ctx, plan.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.CreateConnectionInput{
		Name:        plan.Name.ValueString(),
		Provider:    plan.ProviderName.ValueString(),
		Credentials: creds,
		Config:      config,
	}

	conn, err := r.client.CreateConnection(ctx, plan.ProjectID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating connection", err.Error())
		return
	}

	plan.ID = types.StringValue(conn.ID)
	plan.ProjectID = types.StringValue(conn.ProjectID)
	plan.Name = types.StringValue(conn.Name)
	plan.ProviderName = types.StringValue(conn.Provider)

	configMap, diags := helpers.FlattenStringMap(conn.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Config = configMap

	plan.CreatedAt = types.StringValue(conn.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(conn.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *connectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state connectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn, err := r.client.GetConnection(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if showbiz.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading connection", err.Error())
		return
	}

	state.ID = types.StringValue(conn.ID)
	state.ProjectID = types.StringValue(conn.ProjectID)
	state.Name = types.StringValue(conn.Name)
	state.ProviderName = types.StringValue(conn.Provider)

	configMap, diags := helpers.FlattenStringMap(conn.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Config = configMap

	// Credentials are write-only; preserve state value (API does not return them).
	state.CreatedAt = types.StringValue(conn.CreatedAt.String())
	state.UpdatedAt = types.StringValue(conn.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *connectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan connectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state connectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, diags := helpers.ExpandStringMap(ctx, plan.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.UpdateConnectionInput{
		Config: config,
	}

	conn, err := r.client.UpdateConnection(ctx, state.ProjectID.ValueString(), state.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating connection", err.Error())
		return
	}

	plan.ID = types.StringValue(conn.ID)
	plan.ProjectID = types.StringValue(conn.ProjectID)
	plan.Name = types.StringValue(conn.Name)
	plan.ProviderName = types.StringValue(conn.Provider)

	configMap, diags := helpers.FlattenStringMap(conn.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Config = configMap

	plan.CreatedAt = types.StringValue(conn.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(conn.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *connectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state connectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConnection(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting connection", err.Error())
		return
	}
}

func (r *connectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: <project_id>/<connection_id>, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
