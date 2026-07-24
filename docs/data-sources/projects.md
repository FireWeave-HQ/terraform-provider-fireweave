---
page_title: "fireweave_projects Data Source - fireweave"
subcategory: ""
description: |-
  Lists all FireWeave projects in the authenticated organisation.
---

# fireweave_projects (Data Source)

Returns every project owned by the organisation associated with the API key.

## Example Usage

```terraform
data "fireweave_projects" "all" {}

output "project_slugs" {
  value = [for p in data.fireweave_projects.all.projects : p.slug]
}
```

## Argument Reference

This data source has no configuration arguments.

### Read-Only

- `projects` (Attributes List) Projects in the organisation.
  - `id` (String) Project identifier.
  - `name` (String) Project name.
  - `slug` (String) Project slug.
  - `status` (String) Project status when returned by the API.
