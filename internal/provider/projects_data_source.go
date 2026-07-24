package provider

import (
	"context"
	"fmt"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*ProjectsDataSource)(nil)

type ProjectsDataSource struct {
	client *client.Client
}

type ProjectsDataSourceModel struct {
	Projects types.List `tfsdk:"projects"`
}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

func (d *ProjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all FireWeave projects in the authenticated organisation.",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"slug": schema.StringAttribute{Computed: true},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Present when returned by the API; may be empty for list views.",
						},
					},
				},
			},
		},
	}
}

func (d *ProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *ProjectsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	projects, err := d.client.ListProjects(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing projects", err.Error())
		return
	}

	attrTypes := map[string]attr.Type{
		"id":     types.StringType,
		"name":   types.StringType,
		"slug":   types.StringType,
		"status": types.StringType,
	}
	elems := make([]attr.Value, 0, len(projects))
	for _, p := range projects {
		obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":     types.StringValue(p.ProjectID),
			"name":   types.StringValue(p.Name),
			"slug":   types.StringValue(p.Slug),
			"status": types.StringValue(p.Status),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		elems = append(elems, obj)
	}

	list, diags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, elems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &ProjectsDataSourceModel{Projects: list})...)
}
