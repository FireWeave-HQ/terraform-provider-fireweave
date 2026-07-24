# Examples

Runnable samples for the FireWeave Terraform provider.

| Path | Description |
|------|-------------|
| [`resources/fireweave_project`](./resources/fireweave_project) | Create a project |
| [`resources/fireweave_environment`](./resources/fireweave_environment) | Project + stage/prod environments |
| [`data-sources`](./data-sources) | Look up projects and environments |

## Running an example

```bash
export FIREWEAVE_API_KEY="fw_org_…"
cd examples/resources/fireweave_environment
terraform init
terraform plan
```

For self-hosted or local FireWeave APIs:

```bash
export FIREWEAVE_ENDPOINT="http://localhost:3001"
```
