terraform {
  backend "http" {
    address        = "https://api.abbey.io/terraform-http-backend"
    lock_address   = "https://api.abbey.io/terraform-http-backend/lock"
    unlock_address = "https://api.abbey.io/terraform-http-backend/unlock"
    lock_method    = "POST"
    unlock_method  = "POST"
  }

  required_providers {
    abbey = {
      source = "abbeylabs/abbey"
      version = "0.2.4"
    }
  }
}

locals {
  # Replace if your abbey email doesn't match your Google User email
  # Example: gcp_member = "your-username@gmail.com"
  google_group_member_email = "{{ .data.system.abbey.identities.abbey.email }}"
  google_customer_id = "{{ .GoogleCustomerId }}" # CHANGEME
}

provider "googleworkspace" {
  customer_id = "${local.google_customer_id}"
}

provider "abbey" {
  # Configuration options
  bearer_auth = var.abbey_token
}

resource "googleworkspace_group" "google_workspace_demo" {
  email       = "{{ .GoogleWorkspaceGroup }}" # CHANGEME
}

resource "abbey_grant_kit" "googleworkspace" {
  name = "GoogleWorkspace"
  description = <<-EOT
    Grants access to Abbey's GoogleWorkspace Group for the Quickstart.
  EOT

  workflow = {
    steps = [
      {
        reviewers = {
          one_of = ["{{ .Reviewer }}"] # CHANGEME
        }
      }
    ]
  }

  policies = [
    { bundle = "{{ .PolicyBundle}}" } # CHANGEME
  ]

  output = {
    # Replace with your own path pointing to where you want your access changes to manifest.
    # Path is an RFC 3986 URI, such as `github://{organization}/{repo}/path/to/file.tf`.
    location = "{{ .AccessOutput }}" # CHANGEME
    append = <<-EOT
      resource "googleworkspace_group_member" "member" {
        group_id = googleworkspace_group.google_workspace_demo.id
        email = "${local.google_group_member_email}"
        role = "MEMBER"
      }
    EOT
  }
}