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

    okta = {
      source  = "okta/okta"
      version = "~> 4.0.2"
    }
  }
}

provider "abbey" {
  # Configuration options
  bearer_auth = var.abbey_token
}

provider "okta" {
  # Configuration options
  org_name  = var.okta_org_name
  base_url  = var.okta_base_url
}

resource "abbey_grant_kit" "okta_group_has_nice_things" {
  name = "Okta_Group__Has_Nice_Things"
  description = <<-EOT
    This resource represents a Okta Group Membership for engineers looking to have nice things.
  EOT

  workflow = {
    steps = [
      {
        reviewers = {
          one_of = ["replace-me@example.com"] # CHANGEME
        }
      }
    ]
  }

  output = {
    # Replace with your own path pointing to where you want your access changes to manifest.
    # Path is an RFC 3986 URI, such as `github://{organization}/{repo}/path/to/file.tf`.
    location = "github://replace-me-with-organization/replace-me-with-repo/access.tf" # CHANGEME
    append = <<-EOT
      resource "okta_user_group_memberships" "has_nice_things__{{ .data.system.abbey.identities.okta.user_id }}" { # {{ .data.system.abbey.identities.abbey.email }}
        user_id = "{{ .data.system.abbey.identities.okta.user_id }}"
        groups = ["${data.okta_group.has_nice_things.id}"]
      }
    EOT
  }
}

resource "abbey_identity" "user_1" {
  abbey_account = "replace-me@example.com"
  source = "okta"
  metadata = jsonencode(
    {
      user_id = "00uReplaceWithOktaUserId" #CHANGEME
    }
  )
}

data "okta_group" "has_nice_things" {
  name = "Has nice things"
}
