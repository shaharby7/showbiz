package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
)

var (
	_ datasource.DataSource              = &providerDataSource{}
	_ datasource.DataSourceWithConfigure = &providerDataSource{}
)

type providerDataSource struct {
	client *showbiz.Client
}

type providerDataSourceModel struct {
	Name          types.String `tfsdk:"name"`
	ResourceTypes types.List   `tfsdk:"resource_types"`
}

func NewProviderDataSource() datasource.DataSource {
	return &providerDataSource{}
}

func (d *providerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider"
}

func (d *providerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Look up a Showbiz cloud provider.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the provider.",
				Required:    true,
			},
			"resource_types": schema.ListAttribute{
				Description: "The resource types supported by this provider.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *providerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *providerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p, err := d.client.GetProvider(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading provider", err.Error())
		return
	}

	config.Name = types.StringValue(p.Name)

	rtList, diags := types.ListValueFrom(ctx, types.StringType, p.ResourceTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.ResourceTypes = rtList

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
