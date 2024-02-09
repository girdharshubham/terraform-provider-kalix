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
	_ datasource.DataSource              = &kalixProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &kalixProjectsDataSource{}
)

type kalixProjectsDataSource struct {
	KalixProviderConfig
}

type kalixProjectsDataSourceModel struct {
	Projects types.String `tfsdk:"projects"`
}

//type kalixProjectDataSourceModel struct {
//	ID          types.String `tfsdk:"id"`
//	Name        types.String `tfsdk:"name"`
//	Description types.String `tfsdk:"description"`
//	Owner       types.String `tfsdk:"owner"`
//}

func NewKalixProjectsDataSource() datasource.DataSource {
	return &kalixProjectsDataSource{}
}

func (d *kalixProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *kalixProjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"projects": schema.StringAttribute{
				Optional: true,
			},
			//"name": schema.StringAttribute{
			//	Optional: true,
			//},
			//"description": schema.StringAttribute{
			//	Optional: true,
			//},
			//"owner": schema.StringAttribute{
			//	Optional: true,
			//},
		},
	}
}

func (d *kalixProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state kalixProjectsDataSourceModel

	projects, err := exec.Command(d.KalixProviderConfig.CliPath, "projects", "list", "-o", "json").Output()
	if err != nil {
		resp.Diagnostics.AddError("Unable to get kalix projects", "")
		return
	}

	state.Projects = types.StringValue(strings.Trim(string(projects[:]), "\n"))

	diagnostics := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}
func (d *kalixProjectsDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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
