package provider

import (
	"context"
	"fmt"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*ProjectDataSource)(nil)

type ProjectDataSource struct {
	client *client.Client
}

type ProjectDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Slug        types.String `tfsdk:"slug"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

func (d *ProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up a FireWeave project by `id` or `slug`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Project identifier.",
			},
			"slug": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Project slug.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project name.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project description.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project status.",
			},
		},
	}
}

func (d *ProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ProjectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsNull() && config.Slug.IsNull() {
		resp.Diagnostics.AddError("Missing lookup key", "Set either `id` or `slug`.")
		return
	}

	var project *client.Project
	var err error
	if !config.ID.IsNull() && !config.ID.IsUnknown() && config.ID.ValueString() != "" {
		project, err = d.client.GetProject(ctx, config.ID.ValueString())
	} else {
		projects, listErr := d.client.ListProjects(ctx)
		if listErr != nil {
			resp.Diagnostics.AddError("Error listing projects", listErr.Error())
			return
		}
		slug := config.Slug.ValueString()
		for i := range projects {
			if projects[i].Slug == slug {
				project = &projects[i]
				break
			}
		}
		if project == nil {
			resp.Diagnostics.AddError("Project not found", fmt.Sprintf("No project with slug %q", slug))
			return
		}
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	config.ID = types.StringValue(project.ProjectID)
	config.Slug = types.StringValue(project.Slug)
	config.Name = types.StringValue(project.Name)
	config.Status = types.StringValue(project.Status)
	if project.Description != nil {
		config.Description = types.StringValue(*project.Description)
	} else {
		config.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
