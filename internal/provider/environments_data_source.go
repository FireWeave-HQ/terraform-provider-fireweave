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

var _ datasource.DataSource = (*EnvironmentsDataSource)(nil)

type EnvironmentsDataSource struct {
	client *client.Client
}

type EnvironmentsDataSourceModel struct {
	ProjectID    types.String `tfsdk:"project_id"`
	Environments types.List   `tfsdk:"environments"`
}

func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{}
}

func (d *EnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (d *EnvironmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists FireWeave environments for a project.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project identifier.",
			},
			"environments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true},
						"slug":         schema.StringAttribute{Computed: true},
						"display_name": schema.StringAttribute{Computed: true},
						"is_default":   schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *EnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config EnvironmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envs, err := d.client.ListEnvironments(ctx, config.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error listing environments", err.Error())
		return
	}

	attrTypes := map[string]attr.Type{
		"id":           types.StringType,
		"slug":         types.StringType,
		"display_name": types.StringType,
		"is_default":   types.BoolType,
	}
	elems := make([]attr.Value, 0, len(envs))
	for _, e := range envs {
		obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":           types.StringValue(e.EnvID),
			"slug":         types.StringValue(e.Slug),
			"display_name": types.StringValue(e.DisplayName),
			"is_default":   types.BoolValue(e.IsDefault),
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

	config.Environments = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
