terraform {
  backend "http" {
    address        = "http://localhost:8080/terraform-http-backend"
    lock_address   = "http://localhost:8080/terraform-http-backend/lock"
    unlock_address = "http://localhost:8080/terraform-http-backend/unlock"
    # address        = "https://api.abbey.io/terraform-http-backend"
    # lock_address   = "https://api.abbey.io/terraform-http-backend/lock"
    # unlock_address = "https://api.abbey.io/terraform-http-backend/unlock"
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
  google_group_member_email = "apple@arvil.co"
  google_customer_id = "C022vh2jg" # using arvil.co
}

provider "googleworkspace" {
  customer_id = "${local.google_customer_id}"
}

provider "abbey" {
  # Configuration options
  bearer_auth = var.abbey_token
  server_url = "http://localhost:8080/v1/"
}

# resource "googleworkspace_group" "google_workspace_demo" {
#   email       = "google-workspace-groups-demo@arvil.co"
# }

# data "googleworkspace_group" "abbey_gke_demo" {
#   email = "gke-quickstart@abbey.io"
# }

# output "abbey_gke_group_name" {
#   value = data.googleworkspace_group.abbey_gke_demo.name
# }

# resource "abbey_grant_kit" "googleworkspace" {
#   name = "GoogleWorkspace"
#   description = <<-EOT
#     Grants access to Abbey's GoogleWorkspace Group for the Quickstart.
#   EOT

#   workflow = {
#     steps = [
#       {
#         reviewers = {
#           one_of = ["hat@abbey.io"]
#         }
#       }
#     ]
#   }

#   policies = [
#     { bundle = "github://hatim-khan/abbey-example-google-workspace-groups/policies" }
#   ]

#   output = {
#     # Replace with your own path pointing to where you want your access changes to manifest.
#     # Path is an RFC 3986 URI, such as `github://{organization}/{repo}/path/to/file.tf`.
#     location = "github://hatim-khan/abbey-example-google-workspace-groups/access.tf"
#     append = <<-EOT
#       resource "googleworkspace_group_member" "member" {
#         group_id = googleworkspace_group.google_workspace_demo.id
#         email = "${local.google_group_member_email}"
#         role = "MEMBER"
#       }
#     EOT
#   }
# }

resource "abbey_grant_kit" "abbey_slack" {
  name = "Abbey-Slack"
  description = <<-EOT
    Grants access to Abbey's slack.
  EOT

  workflow = {
    steps = [
      {
        reviewers = {
          one_of = ["hat+test@abbey.io"]
        }
      }
    ]
  }

  policies = [
    { bundle = "github://hatim-khan/abbey-example-google-workspace-groups/policies" }
  ]

  output = {
    # Replace with your own path pointing to where you want your access changes to manifest.
    # Path is an RFC 3986 URI, such as `github://{organization}/{repo}/path/to/file.tf`.
    location = "github://hatim-khan/abbey-example-google-workspace-groups/access.tf"
    append = <<-EOT
      resource "abbey_demo" "grant_read_write_access" {
        permission = "read_write"
        email = "{{ .data.system.abbey.identities.abbey.email }}"
      }
    EOT
  }
}
