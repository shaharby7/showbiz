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
)

var (
	_ resource.Resource                = &iamPolicyResource{}
	_ resource.ResourceWithConfigure   = &iamPolicyResource{}
	_ resource.ResourceWithImportState = &iamPolicyResource{}
)

type iamPolicyResource struct {
	client *showbiz.Client
}

type iamPolicyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Permissions    types.List   `tfsdk:"permissions"`
	Scope          types.String `tfsdk:"scope"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func NewIAMPolicyResource() resource.Resource {
	return &iamPolicyResource{}
}

func (r *iamPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policy"
}

func (r *iamPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Showbiz IAM policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID this policy belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permissions": schema.ListAttribute{
				Description: "The list of permissions granted by this policy.",
				Required:    true,
				ElementType: types.StringType,
			},
			"scope": schema.StringAttribute{
				Description: "The scope of the policy.",
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

func (r *iamPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iamPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	resp.Diagnostics.Append(plan.Permissions.ElementsAs(ctx, &permissions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.CreatePolicyInput{
		Name:        plan.Name.ValueString(),
		Permissions: permissions,
	}

	policy, err := r.client.CreateOrgPolicy(ctx, plan.OrganizationID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM policy", err.Error())
		return
	}

	plan.ID = types.StringValue(policy.ID)
	plan.Name = types.StringValue(policy.Name)
	plan.OrganizationID = types.StringValue(policy.OrganizationID)
	plan.Scope = types.StringValue(policy.Scope)

	permList, diags := types.ListValueFrom(ctx, types.StringType, policy.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Permissions = permList

	plan.CreatedAt = types.StringValue(policy.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(policy.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iamPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iamPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetPolicy(ctx, state.ID.ValueString())
	if err != nil {
		if showbiz.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading IAM policy", err.Error())
		return
	}

	state.ID = types.StringValue(policy.ID)
	state.Name = types.StringValue(policy.Name)
	state.OrganizationID = types.StringValue(policy.OrganizationID)
	state.Scope = types.StringValue(policy.Scope)

	permList, diags := types.ListValueFrom(ctx, types.StringType, policy.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Permissions = permList

	state.CreatedAt = types.StringValue(policy.CreatedAt.String())
	state.UpdatedAt = types.StringValue(policy.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *iamPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iamPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state iamPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	resp.Diagnostics.Append(plan.Permissions.ElementsAs(ctx, &permissions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.UpdatePolicyInput{
		Permissions: permissions,
	}

	policy, err := r.client.UpdateOrgPolicy(ctx, state.OrganizationID.ValueString(), state.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating IAM policy", err.Error())
		return
	}

	plan.ID = types.StringValue(policy.ID)
	plan.Name = types.StringValue(policy.Name)
	plan.OrganizationID = types.StringValue(policy.OrganizationID)
	plan.Scope = types.StringValue(policy.Scope)

	permList, diags := types.ListValueFrom(ctx, types.StringType, policy.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Permissions = permList

	plan.CreatedAt = types.StringValue(policy.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(policy.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iamPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iamPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgPolicy(ctx, state.OrganizationID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting IAM policy", err.Error())
		return
	}
}

func (r *iamPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: <organization_id>/<policy_id>, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
