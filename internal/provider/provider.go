package provider

import (
	"context"
	"os"

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
	Path types.String `tfsdk:"path"`
}

func (k *kalixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "kalix"
	response.Version = k.version
}

func (k *kalixProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
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

	if kpm.Path.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Unknown path for the kalix cli",
			"use the KALIX_PATH environment variable.")
	}

	if response.Diagnostics.HasError() {
		return
	}
	kalixPath := os.Getenv("KALIX_PATH")
	if !kpm.Path.IsNull() {
		kalixPath = kpm.Path.ValueString()
	}

	if kalixPath == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Missing Kalix CLI Path", "",
		)
	}

	response.DataSourceData = kalixPath
	response.ResourceData = kalixPath

}

func (k *kalixProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (k *kalixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKalixCliVersionDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kalixProvider{
			version: version,
		}
	}
}
