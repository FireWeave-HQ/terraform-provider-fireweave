---
page_title: "FireWeave Provider"
description: |-
  Manage FireWeave projects and environments.
---

# FireWeave Provider

The FireWeave provider manages projects and environments via the FireWeave control-plane `/v1` API using an org-scoped `fw_org_…` API key.

## Example Usage

```terraform
provider "fireweave" {
  # endpoint = "https://app-server.fireweave.ai"
  # api_key  = var.fireweave_api_key
}
```

## Schema

### Optional

- `api_key` (String, Sensitive) Org-scoped FireWeave API key (`fw_org_…`). Can also be set via `FIREWEAVE_API_KEY`.
- `endpoint` (String) FireWeave API endpoint (default `https://app-server.fireweave.ai`). Can also be set via `FIREWEAVE_ENDPOINT`.
