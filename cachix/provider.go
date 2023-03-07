package cachix

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	apiclient "github.com/autophagy/terraform-provider-cachix/client"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type cachixProvider struct{}

var (
	_ provider.Provider = &cachixProvider{}
)

type cachixProviderModel struct {
	ApiToken types.String `tfsdk:"apiToken"`
}

func (p *cachixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cachix"
}

func (p *cachixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cachix",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "Cachix API Token. Can also be provided via $CACHIX_API_TOKEN.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *cachixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config cachixProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := os.Getenv("CACHIX_API_TOKEN")

	if !config.ApiToken.IsNull() {
		apiToken = config.ApiToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Cachix API Token",
			"Either set the value or use the CACHIX_API_TOKEN env var.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cachix_api_token", apiToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cachix_api_token")

	transport := httptransport.New("api.cachix.org", "", nil)
	client := apiclient.New(transport, strfmt.Default)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *cachixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *cachixProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func New() provider.Provider {
	return &cachixProvider{}
}
