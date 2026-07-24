terraform {
  required_providers {
    fireweave = {
      source = "FireWeave-HQ/fireweave"
    }
  }
}

provider "fireweave" {
  # endpoint = "https://app-server.fireweave.ai" # optional
  # api_key  = var.fireweave_api_key            # or FIREWEAVE_API_KEY
}

resource "fireweave_project" "demo" {
  name        = "Demo Project"
  slug        = "demo-project"
  description = "Managed by Terraform"
}
