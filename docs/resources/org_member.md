---
page_title: "fireweave_org_member Resource - fireweave"
subcategory: ""
description: |-
  Manages an organisation member role assignment and removal.
---

# fireweave_org_member (Resource)

Manages an organisation **member role assignment**.

Create and update call `PUT /api/organizations/:orgId/members/:userId/role`. Destroy calls `DELETE /api/organizations/:orgId/members/:userId`.

The target user must already be a member of the organisation (for example after accepting an invite). This resource does not send invitations. The organisation `owner` role cannot be assigned or changed here.

## Example Usage

```terraform
resource "fireweave_org_member" "alice_admin" {
  organization_id = "01234567-89ab-cdef-0123-456789abcdef"
  user_id         = "11111111-2222-3333-4444-555555555555"
  role            = "admin"
}
```

## Argument Reference

### Required

- `organization_id` (String) Organisation identifier. Changing this forces a new resource.
- `user_id` (String) User identifier of the member. Changing this forces a new resource.
- `role` (String) Role to assign. One of `admin`, `member`, or `viewer`.

### Read-Only

- `id` (String) Terraform identifier (`organization_id/user_id`).
- `member_id` (String) Membership row identifier returned by FireWeave.

## Import

Import by organisation id and user id:

```shell
terraform import fireweave_org_member.alice_admin 01234567-89ab-cdef-0123-456789abcdef/11111111-2222-3333-4444-555555555555
```
