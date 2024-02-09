package provider

import (
	"context"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &kalixProvider{}
)

type kalixProvider struct {
	version string
}

type kalixProviderModel struct {
	CliPath      types.String `tfsdk:"cli_path"`
	RefreshToken types.String `tfsdk:"refresh_token"`
}

func (k *kalixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "kalix"
	response.Version = k.version
}

func (k *kalixProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cli_path": schema.StringAttribute{
				Optional: true,
			},
			"refresh_token": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (k *kalixProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var kpm kalixProviderModel
	diagnostics := request.Config.Get(ctx, &kpm)
	response.Diagnostics.Append(diagnostics...)
	if response.Diagnostics.HasError() {
		return
	}

	if kpm.CliPath.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("cli_path"), "", "")
	}

	if kpm.RefreshToken.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("refresh_token"), "", "")
	}

	if response.Diagnostics.HasError() {
		return
	}

	kalix, _ := exec.LookPath("kalix")
	if !kpm.CliPath.IsNull() {
		kalix = kpm.CliPath.ValueString()
	}

	if kalix == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("cli_path"),
			"Missing Kalix CLI Path",
			"",
		)
	}

	refreshToken := os.Getenv("KALIX_REFRESH_TOKEN")
	if !kpm.RefreshToken.IsNull() {
		refreshToken = kpm.RefreshToken.ValueString()
	}

	if refreshToken == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Missing Kalix Refresh Token",
			"")
	}

	kalicProviderConfig := KalixProviderConfig{
		CliPath:      kalix,
		RefreshToken: refreshToken,
	}

	response.DataSourceData = kalicProviderConfig
	response.ResourceData = kalicProviderConfig

}

func (k *kalixProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (k *kalixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKalixCliVersionDataSource,
		NewKalixProjectsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kalixProvider{
			version: version,
		}
	}
}
