terraform {
  required_version = ">= 1.5.0"

  required_providers {
    fireweave = {
      source  = "FireWeave-HQ/fireweave"
      version = "~> 0.1"
    }
  }
}

provider "fireweave" {}

resource "fireweave_project" "demo" {
  name        = "Demo Project"
  slug        = "demo-project"
  description = "Promotion pipeline managed by Terraform"
}

resource "fireweave_environment" "stage" {
  project_id   = fireweave_project.demo.id
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

  webhook_aliases = ["staging"]
}

resource "fireweave_environment" "prod" {
  project_id   = fireweave_project.demo.id
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

output "stage_env_id" {
  value = fireweave_environment.stage.id
}

output "prod_env_id" {
  value = fireweave_environment.prod.id
}
