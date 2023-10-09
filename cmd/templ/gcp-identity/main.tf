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

provider "google" {
  billing_project     = "replace-me"
  region      = "us-west1"
}

provider "abbey" {
  # Configuration options
  bearer_auth = var.abbey_token
}

locals {
  # Replace if your abbey email doesn't match your Google User email
  # Example: gcp_member = "your-username@gmail.com"
  gcp_member = "{{ .data.system.abbey.identities.abbey.email }}"
  gcp_customer_id = "replace-me"
}


resource "google_cloud_identity_group" "abbey-gcp-quickstart" {
  display_name         = "abbey-gcp-quickstart"
  initial_group_config = "WITH_INITIAL_OWNER"

  parent = "customers/${local.gcp_customer_id}"

  group_key {
    # choose a unique group ID
    id = "replace-me@example.com"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "abbey_grant_kit" "abbey_emergency" {
  name = "Abbey-Emergency-Access"
  description = <<-EOT
    Grants access to Abbey's GCP Group for emergency scenarios.
  EOT

  workflow = {
    steps = [
      {
        reviewers = {
          one_of = ["{{ .Reviewer }}"]
        }
      }
    ]
  }

  policies = [
    { bundle = "{{ .PolicyBundle}}" }
  ]

  output = {
    # Replace with your own path pointing to where you want your access changes to manifest.
    # Path is an RFC 3986 URI, such as `github://{organization}/{repo}/path/to/file.tf`.
    location = "{{ .AccessOutput }}"
    append = <<-EOT
      resource "google_cloud_identity_group_membership" "member" {
        group    = google_cloud_identity_group.abbey_emergency_group.id
        roles {
          name = "MEMBER"
        }
        preferred_member_key {
          id = "${local.gcp_member}"
        }
      }
    EOT
  }
}
