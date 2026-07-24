---
page_title: "fireweave_environment Resource - fireweave"
subcategory: ""
description: |-
  Manages a FireWeave project environment.
---

# fireweave_environment (Resource)

Manages a FireWeave project environment (deploy stage).

## Example Usage

```terraform
resource "fireweave_environment" "stage" {
  project_id   = fireweave_project.demo.id
  slug         = "stage"
  display_name = "Stage"
  is_default   = true

  branch_rules = [
    { kind = "exact", value = "main" }
  ]
}
```

## Schema

### Required

- `project_id` (String) Owning project identifier. Forces replacement.
- `slug` (String) Environment slug (unique within the project). Forces replacement.
- `display_name` (String) Human-readable environment name.

### Optional

- `is_default` (Boolean) Whether this environment is the project default.
- `branch_rules` (Attributes List) Git branch matching rules.
- `tag_rules` (Attributes List) Git tag matching rules.
- `webhook_aliases` (List of String) Optional webhook alias strings.

### Read-Only

- `id` (String) Environment identifier (`envId`).

## Import

```shell
terraform import fireweave_environment.stage <project_id>/<env_id>
```
