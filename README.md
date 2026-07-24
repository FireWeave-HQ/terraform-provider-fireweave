# Terraform Provider for FireWeave

Manage FireWeave **projects** and **environments** with Terraform.

Published as `FireWeave-HQ/fireweave` on the [Terraform Registry](https://registry.terraform.io/).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.5
- A FireWeave org API key (`fw_org_…`), minted via:

```http
POST /api/organizations/:orgId/api-keys
Authorization: <session cookie>
{ "name": "terraform" }
```

The feature is flag-gated (`org-api-keys-v1-management`) on the FireWeave server.

## Example

```hcl
terraform {
  required_providers {
    fireweave = {
      source  = "FireWeave-HQ/fireweave"
      version = ">= 0.1.0"
    }
  }
}

provider "fireweave" {
  # endpoint = "https://app-server.fireweave.ai"
  # api_key  = "fw_org_…"   # or FIREWEAVE_API_KEY
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
    { kind = "exact", value = "main" }
  ]
}
```

## Provider configuration

| Argument   | Env var             | Default                           |
|------------|---------------------|-----------------------------------|
| `endpoint` | `FIREWEAVE_ENDPOINT`| `https://app-server.fireweave.ai` |
| `api_key`  | `FIREWEAVE_API_KEY` | _(required)_                      |

## Resources

- `fireweave_project` — create/update/delete/import by project id
- `fireweave_environment` — create/update/delete/import as `project_id/env_id`

## Data sources

- `fireweave_project` — lookup by `id` or `slug`
- `fireweave_projects` — list all projects in the org
- `fireweave_environments` — list environments for a `project_id`

## Local development

```bash
go build -o terraform-provider-fireweave
# optional: terraform acceptance tests against a local fw-server
export TF_ACC=1
export FIREWEAVE_ENDPOINT=http://localhost:3001
export FIREWEAVE_API_KEY=fw_org_…
go test ./internal/provider -v -timeout 30m
```

## Publishing to the Terraform Registry

These steps are **manual** (require org admin access):

1. **Create the public GitHub repo**  
   `FireWeave-HQ/terraform-provider-fireweave` (name must match `terraform-provider-<NAME>`).

2. **Push this codebase** to that repo (`main` branch).

3. **Generate a GPG signing key** for releases and export the private key:

   ```bash
   gpg --full-generate-key
   gpg --armor --export-secret-keys <KEY_ID>
   ```

4. **Add GitHub Actions secrets** on the repo:
   - `GPG_PRIVATE_KEY` — armored private key
   - `PASSPHRASE` — key passphrase

5. **Publish the GPG public key** to a keyserver (e.g. `keys.openpgp.org`) and note the fingerprint.

6. **Register the provider** at [registry.terraform.io](https://registry.terraform.io/publish/provider):
   - Sign in with the FireWeave-HQ GitHub org
   - Select `terraform-provider-fireweave`
   - Paste the GPG public key / fingerprint

7. **Cut a release tag**:

   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

   The `release` workflow runs GoReleaser, signs the checksums, and creates a GitHub Release. The Registry picks it up automatically once registered.

## License

MPL-2.0
