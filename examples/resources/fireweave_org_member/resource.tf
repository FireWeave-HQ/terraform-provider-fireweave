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

resource "fireweave_org_member" "demo" {
  organization_id = var.organization_id
  user_id         = var.user_id
  role            = "member"
}

variable "organization_id" {
  type        = string
  description = "FireWeave organisation identifier"
}

variable "user_id" {
  type        = string
  description = "Existing organisation member user identifier"
}

output "member_id" {
  description = "Membership row identifier"
  value       = fireweave_org_member.demo.member_id
}
