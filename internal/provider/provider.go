package provider

import (
	"context"
	"os"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultEndpoint = "https://app-server.fireweave.ai"

// Ensure FireweaveProvider satisfies provider.Provider.
var _ provider.Provider = (*FireweaveProvider)(nil)

type FireweaveProvider struct {
	version string
}

type FireweaveProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FireweaveProvider{version: version}
	}
}

func (p *FireweaveProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fireweave"
	resp.Version = p.version
}

func (p *FireweaveProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The FireWeave provider manages projects and environments via the FireWeave control-plane `/v1` API using an org-scoped `fw_org_…` API key.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "FireWeave API endpoint (default `https://app-server.fireweave.ai`). Can also be set via `FIREWEAVE_ENDPOINT`.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Org-scoped FireWeave API key (`fw_org_…`). Can also be set via `FIREWEAVE_API_KEY`.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *FireweaveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config FireweaveProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := defaultEndpoint
	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() && config.Endpoint.ValueString() != "" {
		endpoint = config.Endpoint.ValueString()
	} else if v := os.Getenv("FIREWEAVE_ENDPOINT"); v != "" {
		endpoint = v
	}

	apiKey := ""
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = config.APIKey.ValueString()
	}
	if apiKey == "" {
		apiKey = os.Getenv("FIREWEAVE_API_KEY")
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing FireWeave API key",
			"Set the provider `api_key` argument or the `FIREWEAVE_API_KEY` environment variable.",
		)
		return
	}

	c := client.New(endpoint, apiKey)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *FireweaveProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewEnvironmentResource,
	}
}

func (p *FireweaveProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
		NewEnvironmentsDataSource,
	}
}
