terraform {
  required_version = ">= 1.9.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
    }
  }

  # Don't need to define cloud block when using VCS-driven workflow.
  # cloud {
  #   organization = "nhat_org"
  #   workspaces {
  #     name = "ecommerce-go-app"
  #   }
  # }
}