---
page_title: "fireweave_environment Resource - fireweave"
subcategory: ""
description: |-
  Manages a FireWeave project environment (deploy / promotion stage).
---

# fireweave_environment (Resource)

Manages a FireWeave **environment** — a named stage in a project's promotion pipeline (for example `dev`, `stage`, or `prod`).

Environments can declare branch and tag matching rules so FireWeave can resolve which stage a git ref belongs to.

## Example Usage

```terraform
resource "fireweave_project" "app" {
  name = "Payments Service"
  slug = "payments"
}

resource "fireweave_environment" "stage" {
  project_id   = fireweave_project.app.id
  slug         = "stage"
  display_name = "Stage"
  is_default   = true

  branch_rules = [
    {
      kind  = "exact"
      value = "main"
    },
  ]

  tag_rules = [
    {
      kind  = "prefix"
      value = "stage-"
    },
  ]

  webhook_aliases = ["staging", "stg"]
}

resource "fireweave_environment" "prod" {
  project_id   = fireweave_project.app.id
  slug         = "prod"
  display_name = "Production"

  branch_rules = [
    {
      kind  = "exact"
      value = "release"
    },
  ]

  tag_rules = [
    {
      kind  = "glob"
      value = "v*"
    },
  ]
}
```

## Argument Reference

### Required

- `project_id` (String) ID of the owning project. Changing this forces a new resource.
- `slug` (String) Environment slug, unique within the project (for example `stage`). Changing this forces a new resource.
- `display_name` (String) Human-readable name shown in the FireWeave UI.

### Optional

- `is_default` (Boolean) Whether this environment is the project's default. Defaults to `false`.
- `branch_rules` (Attributes List) Rules that map git branches to this environment. See [Ref Rule](#ref-rule) below.
- `tag_rules` (Attributes List) Rules that map git tags to this environment. See [Ref Rule](#ref-rule) below.
- `webhook_aliases` (List of String) Alternate names accepted from deploy / webhook payloads.

### Read-Only

- `id` (String) Environment identifier (`envId`) assigned by FireWeave.

### Ref Rule

Nested object used by `branch_rules` and `tag_rules`:

- `kind` (String) Match strategy. One of:
  - `exact` — full string match
  - `prefix` — value is a prefix of the ref name
  - `glob` — shell-style glob
- `value` (String) Pattern to match.
- `repo` (String, Optional) Limit the rule to a specific repository full name (for example `acme/payments`).

## Import

Import with `<project_id>/<env_id>`:

```shell
terraform import fireweave_environment.stage 01234567-89ab-cdef-0123-456789abcdef/env_stage01
```
