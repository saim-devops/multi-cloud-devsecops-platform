terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.50"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.30"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.13"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }

  # Production remote state with S3 and DynamoDB locking.
  # Uncomment when remote bucket is created:
  # backend "s3" {
  #   bucket         = "enterprise-devsecops-tfstate-aws"
  #   key            = "environments/aws/terraform.tfstate"
  #   region         = "us-east-1"
  #   dynamodb_table = "enterprise-devsecops-tflocks-aws"
  #   encrypt        = true
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "MultiCloudDevSecOps"
      ManagedBy   = "Terraform"
    }
  }
}
