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
	_ datasource.DataSource              = &connectionDataSource{}
	_ datasource.DataSourceWithConfigure = &connectionDataSource{}
)

type connectionDataSource struct {
	client *showbiz.Client
}

type connectionDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Name         types.String `tfsdk:"name"`
	ProviderName types.String `tfsdk:"provider_name"`
	Config       types.Map    `tfsdk:"config"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func NewConnectionDataSource() datasource.DataSource {
	return &connectionDataSource{}
}

func (d *connectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (d *connectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Look up a Showbiz connection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the connection.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID this connection belongs to.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the connection.",
				Computed:    true,
			},
			"provider_name": schema.StringAttribute{
				Description: "The cloud provider type.",
				Computed:    true,
			},
			"config": schema.MapAttribute{
				Description: "Provider configuration.",
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

func (d *connectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *connectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config connectionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn, err := d.client.GetConnection(ctx, config.ProjectID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading connection", err.Error())
		return
	}

	config.ID = types.StringValue(conn.ID)
	config.ProjectID = types.StringValue(conn.ProjectID)
	config.Name = types.StringValue(conn.Name)
	config.ProviderName = types.StringValue(conn.Provider)

	configMap, diags := helpers.FlattenStringMap(conn.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Config = configMap

	config.CreatedAt = types.StringValue(conn.CreatedAt.String())
	config.UpdatedAt = types.StringValue(conn.UpdatedAt.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
