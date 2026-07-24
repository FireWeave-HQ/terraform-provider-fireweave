---
page_title: "fireweave_project Resource - fireweave"
subcategory: ""
description: |-
  Manages a FireWeave project.
---

# fireweave_project (Resource)

Manages a FireWeave project.

## Example Usage

```terraform
resource "fireweave_project" "demo" {
  name        = "Demo Project"
  slug        = "demo-project"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Human-readable project name.
- `slug` (String) URL-safe project slug (unique within the organisation). Forces replacement.

### Optional

- `description` (String) Optional project description.
- `status` (String) Project status: `active` or `archived`.

### Read-Only

- `id` (String) Project identifier.

## Import

```shell
terraform import fireweave_project.demo <project_id>
```
