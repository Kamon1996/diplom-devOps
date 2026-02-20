terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # --- S3 Backend (optional) ---
  # Uncomment after creating the bucket:
  #   aws s3api create-bucket --bucket habits-tracker-tfstate --region us-east-1
  #   aws s3api put-bucket-versioning --bucket habits-tracker-tfstate \
  #     --versioning-configuration Status=Enabled
  #
  # backend "s3" {
  #   bucket = "habits-tracker-tfstate"
  #   key    = "infrastructure/terraform.tfstate"
  #   region = "us-east-1"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project   = var.project_name
      ManagedBy = "terraform"
    }
  }
}
