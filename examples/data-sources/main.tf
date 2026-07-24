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

# List every project in the organisation.
data "fireweave_projects" "all" {}

# Look up one project by slug (set to a project that already exists).
data "fireweave_project" "selected" {
  slug = var.project_slug
}

data "fireweave_environments" "pipeline" {
  project_id = data.fireweave_project.selected.id
}

variable "project_slug" {
  type        = string
  description = "Slug of an existing FireWeave project to inspect"
}

output "all_project_slugs" {
  value = [for p in data.fireweave_projects.all.projects : p.slug]
}

output "selected_project" {
  value = {
    id     = data.fireweave_project.selected.id
    name   = data.fireweave_project.selected.name
    status = data.fireweave_project.selected.status
  }
}

output "environment_slugs" {
  value = [for e in data.fireweave_environments.pipeline.environments : e.slug]
}
