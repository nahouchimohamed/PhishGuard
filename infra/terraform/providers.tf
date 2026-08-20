provider "aws" {
  region                      = var.aws_region
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  default_tags {
    tags = {
      Project   = "phishguard"
      Environment = "localstack"
      ManagedBy = "terraform"
    }
  }

  endpoints {
    sqs    = var.aws_endpoint
    lambda = var.aws_endpoint
    iam    = var.aws_endpoint
    sts    = var.aws_endpoint
  }
}
