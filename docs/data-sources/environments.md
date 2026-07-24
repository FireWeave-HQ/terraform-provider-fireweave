---
page_title: "fireweave_environments Data Source - fireweave"
subcategory: ""
description: |-
  Lists FireWeave environments for a project.
---

# fireweave_environments (Data Source)

Lists deploy environments for a given project.

## Example Usage

```terraform
data "fireweave_project" "app" {
  slug = "payments"
}

data "fireweave_environments" "pipeline" {
  project_id = data.fireweave_project.app.id
}

output "default_env" {
  value = one([
    for e in data.fireweave_environments.pipeline.environments : e.slug
    if e.is_default
  ])
}
```

## Argument Reference

### Required

- `project_id` (String) Project whose environments should be listed.

### Read-Only

- `environments` (Attributes List) Environments in promotion order (as returned by the API).
  - `id` (String) Environment identifier when available from the API.
  - `slug` (String) Environment slug.
  - `display_name` (String) Human-readable name.
  - `is_default` (Boolean) Whether the environment is the project default.
