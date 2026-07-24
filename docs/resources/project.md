---
page_title: "fireweave_project Resource - fireweave"
subcategory: ""
description: |-
  Manages a FireWeave project.
---

# fireweave_project (Resource)

Manages a FireWeave project within the authenticated organisation.

A project is the top-level container for environments, integrations, and rollouts.

## Example Usage

```terraform
resource "fireweave_project" "app" {
  name        = "Payments Service"
  slug        = "payments"
  description = "Checkout and billing services"
  status      = "active"
}
```

## Argument Reference

### Required

- `name` (String) Human-readable project name.
- `slug` (String) URL-safe slug, unique within the organisation. Changing this forces a new resource.

### Optional

- `description` (String) Longer description of the project.
- `status` (String) Lifecycle status. One of `active` or `archived`. Defaults to `active`.

### Read-Only

- `id` (String) Project identifier assigned by FireWeave.

## Import

Import by project id:

```shell
terraform import fireweave_project.app 01234567-89ab-cdef-0123-456789abcdef
```
