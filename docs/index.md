---
page_title: "FireWeave Provider"
description: |-
  The FireWeave provider manages projects and environments in FireWeave via org-scoped API keys.
---

# FireWeave Provider

The FireWeave provider lets you manage [FireWeave](https://fireweave.ai) **projects** and **environments** with Terraform.

Use it to define promotion pipelines (for example `dev` → `stage` → `prod`) as code, keep environment branch/tag rules in sync across teams, and import existing projects into Terraform state.

## Example Usage

```terraform
terraform {
  required_providers {
    fireweave = {
      source  = "FireWeave-HQ/fireweave"
      version = "~> 0.1"
    }
  }
}

provider "fireweave" {
  # api_key  = var.fireweave_api_key # optional if FIREWEAVE_API_KEY is set
  # endpoint = "https://app-server.fireweave.ai"
}

resource "fireweave_project" "app" {
  name        = "My App"
  slug        = "my-app"
  description = "Managed by Terraform"
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
}
```

## Authentication

Authenticate with an org-scoped API key (`fw_org_…`). Prefer the `FIREWEAVE_API_KEY` environment variable over committing secrets to configuration.

### Creating an API key

1. Sign in to FireWeave for the target organisation.
2. Call:

```http
POST /api/organizations/{orgId}/api-keys
Content-Type: application/json

{ "name": "terraform" }
```

3. Persist the returned `key` securely. It is only returned once.

Org API key management requires the FireWeave `org-api-keys-v1-management` feature to be enabled for the organisation.

## Schema

### Optional

- `api_key` (String, Sensitive) Org-scoped FireWeave API key (`fw_org_…`). Defaults to the `FIREWEAVE_API_KEY` environment variable.
- `endpoint` (String) FireWeave API base URL. Defaults to `https://app-server.fireweave.ai`, or `FIREWEAVE_ENDPOINT` when set.
