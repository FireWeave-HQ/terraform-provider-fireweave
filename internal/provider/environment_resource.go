package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = (*EnvironmentResource)(nil)
	_ resource.ResourceWithImportState = (*EnvironmentResource)(nil)
)

type EnvironmentResource struct {
	client *client.Client
}

type EnvironmentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ProjectID      types.String `tfsdk:"project_id"`
	Slug           types.String `tfsdk:"slug"`
	DisplayName    types.String `tfsdk:"display_name"`
	IsDefault      types.Bool   `tfsdk:"is_default"`
	BranchRules    types.List   `tfsdk:"branch_rules"`
	TagRules       types.List   `tfsdk:"tag_rules"`
	WebhookAliases types.List   `tfsdk:"webhook_aliases"`
}

func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{}
}

func refRuleNestedBlock() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Git ref matching rules (`exact`, `prefix`, or `glob`).",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"kind": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Rule kind: `exact`, `prefix`, or `glob`.",
				},
				"value": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Rule value.",
				},
				"repo": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Optional repo scope for the rule.",
				},
			},
		},
	}
}

func (r *EnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *EnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a FireWeave project environment (deploy stage).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Environment identifier (`envId`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Owning project identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slug": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Environment slug (unique within the project).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable environment name.",
			},
			"is_default": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether this environment is the project default.",
			},
			"branch_rules":    refRuleNestedBlock(),
			"tag_rules":       refRuleNestedBlock(),
			"webhook_aliases": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Optional webhook alias strings for this environment.",
			},
		},
	}
}

func (r *EnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvironmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.CreateEnvironmentInput{
		Slug:        plan.Slug.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
	}
	if !plan.IsDefault.IsNull() && !plan.IsDefault.IsUnknown() {
		v := plan.IsDefault.ValueBool()
		in.IsDefault = &v
	}
	branchRules, diags := listToRefRules(ctx, plan.BranchRules)
	resp.Diagnostics.Append(diags...)
	tagRules, diags := listToRefRules(ctx, plan.TagRules)
	resp.Diagnostics.Append(diags...)
	aliases, diags := listToStrings(ctx, plan.WebhookAliases)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	in.BranchRules = branchRules
	in.TagRules = tagRules
	in.WebhookAliases = aliases

	created, err := r.client.CreateEnvironment(ctx, plan.ProjectID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating environment", err.Error())
		return
	}

	resp.Diagnostics.Append(setEnvironmentModel(ctx, &plan, created)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.client.GetEnvironment(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	resp.Diagnostics.Append(setEnvironmentModel(ctx, &state, env)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvironmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	displayName := plan.DisplayName.ValueString()
	in := client.UpdateEnvironmentInput{DisplayName: &displayName}
	if !plan.IsDefault.IsNull() && !plan.IsDefault.IsUnknown() {
		v := plan.IsDefault.ValueBool()
		in.IsDefault = &v
	}
	branchRules, diags := listToRefRules(ctx, plan.BranchRules)
	resp.Diagnostics.Append(diags...)
	tagRules, diags := listToRefRules(ctx, plan.TagRules)
	resp.Diagnostics.Append(diags...)
	aliases, diags := listToStrings(ctx, plan.WebhookAliases)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	in.BranchRules = branchRules
	in.TagRules = tagRules
	in.WebhookAliases = aliases

	updated, err := r.client.UpdateEnvironment(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error updating environment", err.Error())
		return
	}

	resp.Diagnostics.Append(setEnvironmentModel(ctx, &plan, updated)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteEnvironment(ctx, state.ProjectID.ValueString(), state.ID.ValueString()); err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting environment", err.Error())
	}
}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: "<project_id>/<env_id>"
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			`Expected format "<project_id>/<env_id>".`,
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func setEnvironmentModel(ctx context.Context, m *EnvironmentResourceModel, e *client.Environment) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ID = types.StringValue(e.EnvID)
	m.ProjectID = types.StringValue(e.ProjectID)
	m.Slug = types.StringValue(e.Slug)
	m.DisplayName = types.StringValue(e.DisplayName)
	m.IsDefault = types.BoolValue(e.IsDefault)

	branchRules, d := refRulesToList(ctx, e.BranchRules)
	diags.Append(d...)
	m.BranchRules = branchRules

	tagRules, d := refRulesToList(ctx, e.TagRules)
	diags.Append(d...)
	m.TagRules = tagRules

	aliases, d := stringsToList(ctx, e.WebhookAliases)
	diags.Append(d...)
	m.WebhookAliases = aliases

	return diags
}

func listToRefRules(ctx context.Context, list types.List) ([]client.RefRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var elems []attr.Value
	diags.Append(list.ElementsAs(ctx, &elems, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]client.RefRule, 0, len(elems))
	for _, elem := range elems {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		attrs := obj.Attributes()
		rule := client.RefRule{
			Kind:  attrs["kind"].(types.String).ValueString(),
			Value: attrs["value"].(types.String).ValueString(),
		}
		if repo, ok := attrs["repo"].(types.String); ok && !repo.IsNull() && !repo.IsUnknown() {
			v := repo.ValueString()
			rule.Repo = &v
		}
		out = append(out, rule)
	}
	return out, diags
}

func refRulesToList(ctx context.Context, rules []client.RefRule) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"kind":  types.StringType,
		"value": types.StringType,
		"repo":  types.StringType,
	}
	if rules == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), nil
	}
	elems := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		repo := types.StringNull()
		if rule.Repo != nil {
			repo = types.StringValue(*rule.Repo)
		}
		obj, d := types.ObjectValue(attrTypes, map[string]attr.Value{
			"kind":  types.StringValue(rule.Kind),
			"value": types.StringValue(rule.Value),
			"repo":  repo,
		})
		if d.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), d
		}
		elems = append(elems, obj)
	}
	return types.ListValue(types.ObjectType{AttrTypes: attrTypes}, elems)
}

func listToStrings(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var out []string
	diags.Append(list.ElementsAs(ctx, &out, false)...)
	return out, diags
}

func stringsToList(ctx context.Context, values []string) (types.List, diag.Diagnostics) {
	_ = ctx
	if values == nil {
		return types.ListNull(types.StringType), nil
	}
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	return types.ListValue(types.StringType, elems)
}
