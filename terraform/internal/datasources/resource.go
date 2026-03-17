package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"

	"github.com/shaharby7/showbiz/terraform/internal/helpers"
)

var (
	_ datasource.DataSource              = &resourceDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceDataSource{}
)

type resourceDataSource struct {
	client *showbiz.Client
}

type resourceDataSourceModel struct {
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

func NewResourceDataSource() datasource.DataSource {
	return &resourceDataSource{}
}

func (d *resourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (d *resourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Look up a Showbiz resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the resource.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID this resource belongs to.",
				Required:    true,
			},
			"connection_id": schema.StringAttribute{
				Description: "The connection ID used to manage this resource.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the resource.",
				Computed:    true,
			},
			"resource_type": schema.StringAttribute{
				Description: "The type of resource.",
				Computed:    true,
			},
			"values": schema.MapAttribute{
				Description: "Resource configuration values.",
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

func (d *resourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*showbiz.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"Expected *showbiz.Client, got unexpected type.",
		)
		return
	}
	d.client = client
}

func (d *resourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config resourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetResource(ctx, config.ProjectID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	config.ID = types.StringValue(res.ID)
	config.ProjectID = types.StringValue(res.ProjectID)
	config.ConnectionID = types.StringValue(res.ConnectionID)
	config.Name = types.StringValue(res.Name)
	config.ResourceType = types.StringValue(res.ResourceType)
	config.Status = types.StringValue(res.Status)

	valuesMap, diags := helpers.FlattenStringMap(res.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Values = valuesMap

	config.CreatedAt = types.StringValue(res.CreatedAt.String())
	config.UpdatedAt = types.StringValue(res.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
