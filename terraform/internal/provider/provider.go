package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"

	"github.com/shaharby7/showbiz/terraform/internal/datasources"
	"github.com/shaharby7/showbiz/terraform/internal/resources"
)

var _ provider.Provider = &showbizProvider{}

type showbizProvider struct{}

type showbizProviderModel struct {
	APIURL   types.String `tfsdk:"api_url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New() provider.Provider {
	return &showbizProvider{}
}

func (p *showbizProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "showbiz"
}

func (p *showbizProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing Showbiz resources.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Description: "The URL of the Showbiz API. Defaults to https://api.showbiz.dev.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username (email) for Showbiz authentication.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for Showbiz authentication.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *showbizProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config showbizProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := "https://api.showbiz.dev"
	if !config.APIURL.IsNull() && !config.APIURL.IsUnknown() {
		apiURL = config.APIURL.ValueString()
	}

	client := showbiz.NewClient(apiURL)

	_, err := client.Login(ctx, showbiz.LoginInput{
		Email:    config.Username.ValueString(),
		Password: config.Password.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to authenticate with Showbiz API",
			fmt.Sprintf("Login failed: %s", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *showbizProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewProjectResource,
		resources.NewConnectionResource,
		resources.NewResourceResource,
		resources.NewIAMPolicyResource,
		resources.NewPolicyAttachmentResource,
	}
}

func (p *showbizProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewProjectDataSource,
		datasources.NewConnectionDataSource,
		datasources.NewResourceDataSource,
		datasources.NewProviderDataSource,
	}
}
