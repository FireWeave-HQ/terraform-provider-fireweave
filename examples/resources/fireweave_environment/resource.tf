terraform {
  required_providers {
    fireweave = {
      source = "FireWeave-HQ/fireweave"
    }
  }
}

provider "fireweave" {}

resource "fireweave_project" "demo" {
  name = "Demo Project"
  slug = "demo-project"
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
    }
  ]
}
