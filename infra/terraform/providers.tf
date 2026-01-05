# Defaults target LocalStack. Override aws_endpoint only if you know what you are doing.
# This project never applies Terraform against real AWS in CI or Compose.

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "aws_endpoint" {
  type        = string
  default     = "http://localhost:4566"
  description = "LocalStack endpoint. Leave as-is for the free local Lambda lesson."
}

provider "aws" {
  region                      = var.aws_region
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    sqs    = var.aws_endpoint
    lambda = var.aws_endpoint
    iam    = var.aws_endpoint
    sts    = var.aws_endpoint
  }
}
