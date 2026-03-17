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
	_ resource.Resource                = &policyAttachmentResource{}
	_ resource.ResourceWithConfigure   = &policyAttachmentResource{}
	_ resource.ResourceWithImportState = &policyAttachmentResource{}
)

type policyAttachmentResource struct {
	client *showbiz.Client
}

type policyAttachmentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	UserID         types.String `tfsdk:"user_id"`
	PolicyID       types.String `tfsdk:"policy_id"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

func NewPolicyAttachmentResource() resource.Resource {
	return &policyAttachmentResource{}
}

func (r *policyAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_attachment"
}

func (r *policyAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attaches a Showbiz IAM policy to a user on a project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the policy attachment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The email of the user to attach the policy to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_id": schema.StringAttribute{
				Description: "The ID of the policy to attach.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The creation timestamp.",
				Computed:    true,
			},
		},
	}
}

func (r *policyAttachmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.AttachPolicyInput{
		UserEmail: plan.UserID.ValueString(),
		PolicyID:  plan.PolicyID.ValueString(),
	}

	attachment, err := r.client.AttachPolicy(
		ctx,
		plan.OrganizationID.ValueString(),
		plan.ProjectID.ValueString(),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy attachment", err.Error())
		return
	}

	plan.ID = types.StringValue(attachment.ID)
	plan.ProjectID = types.StringValue(attachment.ProjectID)
	plan.UserID = types.StringValue(attachment.UserEmail)
	plan.PolicyID = types.StringValue(attachment.PolicyID)
	plan.CreatedAt = types.StringValue(attachment.CreatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attachments, err := r.client.ListPolicyAttachments(
		ctx,
		state.OrganizationID.ValueString(),
		state.ProjectID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading policy attachments", err.Error())
		return
	}

	var found *showbiz.PolicyAttachment
	for _, a := range attachments {
		if a.ID == state.ID.ValueString() {
			found = a
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(found.ID)
	state.ProjectID = types.StringValue(found.ProjectID)
	state.UserID = types.StringValue(found.UserEmail)
	state.PolicyID = types.StringValue(found.PolicyID)
	state.CreatedAt = types.StringValue(found.CreatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyAttachmentResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes require replacement; Update should never be called.
	resp.Diagnostics.AddError(
		"Update not supported",
		"Policy attachments cannot be updated in-place. All changes require replacement.",
	)
}

func (r *policyAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := showbiz.DetachPolicyInput{
		UserEmail: state.UserID.ValueString(),
		PolicyID:  state.PolicyID.ValueString(),
	}

	err := r.client.DetachPolicy(
		ctx,
		state.OrganizationID.ValueString(),
		state.ProjectID.ValueString(),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting policy attachment", err.Error())
		return
	}
}

func (r *policyAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: <organization_id>/<project_id>/<attachment_id>, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}
