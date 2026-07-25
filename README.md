# Terraform Provider for FireWeave

[![Terraform Registry](https://img.shields.io/badge/terraform-registry-623CE4?logo=terraform&logoColor=white)](https://registry.terraform.io/providers/FireWeave-HQ/fireweave/latest)
[![GitHub release](https://img.shields.io/github/v/release/FireWeave-HQ/terraform-provider-fireweave)](https://github.com/FireWeave-HQ/terraform-provider-fireweave/releases)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](./LICENSE)

Manage [FireWeave](https://fireweave.ai) projects and environments as code.

**Registry:** [`FireWeave-HQ/fireweave`](https://registry.terraform.io/providers/FireWeave-HQ/fireweave/latest)

## Requirements

| Tool | Version |
|------|---------|
| [Terraform](https://developer.hashicorp.com/terraform/downloads) | >= 1.5 |
| FireWeave org API key | `fw_org_…` (see [Authentication](#authentication)) |

## Quick start

```hcl
terraform {
  required_providers {
    fireweave = {
      source  = "FireWeave-HQ/fireweave"
      version = "~> 0.1"
    }
  }
}

provider "fireweave" {
  # Prefer FIREWEAVE_API_KEY in the environment over hardcoding.
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

```bash
export FIREWEAVE_API_KEY="fw_org_…"
terraform init
terraform apply
```

## Authentication

The provider authenticates with an **org-scoped API key** (`fw_org_…`).

| Setting | Environment variable | Default |
|---------|----------------------|---------|
| `api_key` (sensitive) | `FIREWEAVE_API_KEY` | _(required)_ |
| `endpoint` | `FIREWEAVE_ENDPOINT` | `https://app-server.fireweave.ai` |

### Creating an API key

1. Sign in to the FireWeave app for your organisation.
2. Create a key with a session-authenticated request:

```http
POST /api/organizations/{orgId}/api-keys
Content-Type: application/json

{ "name": "terraform" }
```

3. Store the returned `key` value once — it is shown only at creation time.
4. Export it as `FIREWEAVE_API_KEY` (or pass `api_key` in the provider block).

> Org API key management is enabled when the FireWeave `org-api-keys-v1-management` feature flag is on for your organisation.

## Resources

| Resource | Description |
|----------|-------------|
| [`fireweave_project`](docs/resources/project.md) | A FireWeave project |
| [`fireweave_environment`](docs/resources/environment.md) | A deploy environment (promotion stage) in a project |
| [`fireweave_org_member`](docs/resources/org_member.md) | Organisation member role assignment and removal |

### Import

```bash
terraform import fireweave_project.app <project_id>
terraform import fireweave_environment.stage <project_id>/<env_id>
terraform import fireweave_org_member.alice <organization_id>/<user_id>
```

## Data sources

| Data source | Description |
|-------------|-------------|
| [`fireweave_project`](docs/data-sources/project.md) | Look up one project by `id` or `slug` |
| [`fireweave_projects`](docs/data-sources/projects.md) | List all projects in the organisation |
| [`fireweave_environments`](docs/data-sources/environments.md) | List environments for a project |

## Documentation

- [Terraform Registry docs](https://registry.terraform.io/providers/FireWeave-HQ/fireweave/latest/docs)
- [Examples](./examples)

## Developing the provider

```bash
git clone https://github.com/FireWeave-HQ/terraform-provider-fireweave.git
cd terraform-provider-fireweave
go test ./...
go build -o terraform-provider-fireweave .
```

Acceptance tests (requires a running FireWeave API and a test key):

```bash
export TF_ACC=1
export FIREWEAVE_ENDPOINT=http://localhost:3001
export FIREWEAVE_API_KEY=fw_org_…
make testacc
```

### Cutting a release

Maintainers: tag a semver release (`vX.Y.Z`). The GitHub Actions `release` workflow builds multi-platform binaries with GoReleaser, signs `SHA256SUMS` with the repo GPG key, and publishes a GitHub Release. The Terraform Registry picks up new versions automatically.

Signing material for operators lives under [`.release/`](./.release).

## Support

- Issues: [GitHub Issues](https://github.com/FireWeave-HQ/terraform-provider-fireweave/issues)
- FireWeave product: [fireweave.ai](https://fireweave.ai)

## License

[Mozilla Public License 2.0](./LICENSE)
