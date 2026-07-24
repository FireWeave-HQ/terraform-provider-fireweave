terraform {
  required_version = ">= 1.5.0"

  required_providers {
    fireweave = {
      source  = "FireWeave-HQ/fireweave"
      version = "~> 0.1"
    }
  }
}

provider "fireweave" {
  # Configure via environment variables:
  #   FIREWEAVE_API_KEY      (required)
  #   FIREWEAVE_ENDPOINT     (optional)
}

resource "fireweave_project" "demo" {
  name        = "Demo Project"
  slug        = "demo-project"
  description = "Managed by Terraform"
  status      = "active"
}

output "project_id" {
  description = "FireWeave project identifier"
  value       = fireweave_project.demo.id
}
