package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = (*OrgMemberResource)(nil)
	_ resource.ResourceWithImportState = (*OrgMemberResource)(nil)
)

var validOrgMemberRoles = map[string]struct{}{
	"admin":  {},
	"member": {},
	"viewer": {},
}

type OrgMemberResource struct {
	client *client.Client
}

type OrgMemberResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	UserID         types.String `tfsdk:"user_id"`
	Role           types.String `tfsdk:"role"`
	MemberID       types.String `tfsdk:"member_id"`
}

func NewOrgMemberResource() resource.Resource {
	return &OrgMemberResource{}
}

func (r *OrgMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_member"
}

func (r *OrgMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an organisation member's role assignment. The user must already be a member of the organisation (for example via invite/accept); create and update call `PUT /api/organizations/:orgId/members/:userId/role`, and delete removes the member.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Terraform identifier (`organization_id/user_id`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Organisation identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User identifier of the organisation member.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Member role. One of `admin`, `member`, or `viewer`. The organisation `owner` role cannot be assigned or changed through this resource.",
			},
			"member_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Membership row identifier returned by FireWeave.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *OrgMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrgMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrgMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := plan.Role.ValueString()
	if err := validateOrgMemberRole(role); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("role"), "Invalid role", err.Error())
		return
	}

	orgID := plan.OrganizationID.ValueString()
	userID := plan.UserID.ValueString()
	result, err := r.client.SetOrgMemberRole(ctx, orgID, userID, client.SetOrgMemberRoleInput{Role: role})
	if err != nil {
		resp.Diagnostics.AddError("Error assigning organisation member role", err.Error())
		return
	}

	setOrgMemberModel(&plan, orgID, result.UserID, result.Role, result.MemberID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrgMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrganizationID.ValueString()
	userID := state.UserID.ValueString()
	member, err := r.client.GetOrgMember(ctx, orgID, userID)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading organisation member", err.Error())
		return
	}

	setOrgMemberModel(&state, orgID, member.UserID, member.Role, member.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrgMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrgMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := plan.Role.ValueString()
	if err := validateOrgMemberRole(role); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("role"), "Invalid role", err.Error())
		return
	}

	orgID := plan.OrganizationID.ValueString()
	userID := plan.UserID.ValueString()
	result, err := r.client.SetOrgMemberRole(ctx, orgID, userID, client.SetOrgMemberRoleInput{Role: role})
	if err != nil {
		resp.Diagnostics.AddError("Error updating organisation member role", err.Error())
		return
	}

	setOrgMemberModel(&plan, orgID, result.UserID, result.Role, result.MemberID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrgMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveOrgMember(ctx, state.OrganizationID.ValueString(), state.UserID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error removing organisation member", err.Error())
	}
}

func (r *OrgMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: "<organization_id>/<user_id>"
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			`Expected format "<organization_id>/<user_id>".`,
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func setOrgMemberModel(m *OrgMemberResourceModel, orgID, userID, role, memberID string) {
	m.ID = types.StringValue(orgID + "/" + userID)
	m.OrganizationID = types.StringValue(orgID)
	m.UserID = types.StringValue(userID)
	m.Role = types.StringValue(role)
	m.MemberID = types.StringValue(memberID)
}

func validateOrgMemberRole(role string) error {
	if _, ok := validOrgMemberRoles[role]; !ok {
		return fmt.Errorf("must be one of: admin, member, viewer")
	}
	return nil
}
