---
page_title: "fireweave_project Data Source - fireweave"
subcategory: ""
description: |-
  Looks up a single FireWeave project by id or slug.
---

# fireweave_project (Data Source)

Looks up one FireWeave project. Provide either `id` or `slug`.

## Example Usage

```terraform
data "fireweave_project" "by_slug" {
  slug = "payments"
}

data "fireweave_project" "by_id" {
  id = "01234567-89ab-cdef-0123-456789abcdef"
}

output "project_name" {
  value = data.fireweave_project.by_slug.name
}
```

## Argument Reference

### Optional (one required)

- `id` (String) Project identifier.
- `slug` (String) Project slug within the organisation.

Exactly one of `id` or `slug` should be set.

### Read-Only

- `name` (String) Project name.
- `description` (String) Project description, if set.
- `status` (String) Project status (`active` or `archived`).
- `id` (String) Project identifier (also set when looked up by slug).
- `slug` (String) Project slug (also set when looked up by id).
