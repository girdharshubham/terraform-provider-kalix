// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os/exec"
	"strings"
)

var (
	_ datasource.DataSource              = &kalixCliVersionDataSource{}
	_ datasource.DataSourceWithConfigure = &kalixCliVersionDataSource{}
)

type kalixCliVersionDataSource struct {
	KalixProviderConfig
}

type kalixCliVersionDataSourceModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Version types.String `tfsdk:"version"`
}

func NewKalixCliVersionDataSource() datasource.DataSource {
	return &kalixCliVersionDataSource{}
}

func (d *kalixCliVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cli_version"
}

func (d *kalixCliVersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"version": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (d *kalixCliVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state kalixCliVersionDataSourceModel

	version, err := exec.Command(d.KalixProviderConfig.CliPath, "version").Output()
	if err != nil {
		resp.Diagnostics.AddError("Unable for figure out kalix cli version", "")
		return
	}

	state.Version = types.StringValue(strings.Trim(string(version[:]), "\n"))

	diagnostics := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *kalixCliVersionDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	kalix, ok := request.ProviderData.(KalixProviderConfig)
	if !ok {
		response.Diagnostics.AddError("Unexpected data source configure type", "")
		return
	}

	d.KalixProviderConfig = kalix
}
