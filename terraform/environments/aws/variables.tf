variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS deployment region"
}

variable "environment" {
  type        = string
  default     = "production"
  description = "Environment identifier"
}

variable "cluster_name" {
  type        = string
  default     = "multi-cloud-eks-prod"
  description = "EKS Cluster Name"
}

variable "vpc_cidr" {
  type        = string
  default     = "10.100.0.0/16"
  description = "VPC CIDR block"
}

variable "eks_version" {
  type        = string
  default     = "1.30"
  description = "Kubernetes version for EKS"
}

variable "instance_types" {
  type        = list(string)
  default     = ["t3.medium", "t3a.medium"]
  description = "Node instance types for Spot node group"
}

variable "single_nat_gateway" {
  type        = bool
  default     = true
  description = "Use single NAT Gateway to reduce development cloud costs"
}
