variable "project_id" {
  type        = string
  default     = "enterprise-devsecops-demo"
  description = "GCP Project ID"
}

variable "region" {
  type        = string
  default     = "us-central1"
  description = "Google Cloud region"
}

variable "cluster_name" {
  type        = string
  default     = "multi-cloud-gke-prod"
  description = "GKE Cluster name"
}

variable "network_name" {
  type        = string
  default     = "gke-custom-vpc"
  description = "VPC network name"
}
