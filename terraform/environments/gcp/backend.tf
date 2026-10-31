terraform {
  required_version = ">= 1.5.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.30"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.30"
    }
  }

  # Production remote state backend with GCS:
  # backend "gcs" {
  #   bucket = "enterprise-devsecops-tfstate-gcp"
  #   prefix = "terraform/state/gcp"
  # }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
