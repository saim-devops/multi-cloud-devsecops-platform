variable "location" {
  type        = string
  default     = "eastus"
  description = "Azure Region"
}

variable "resource_group_name" {
  type        = string
  default     = "rg-multi-cloud-aks-prod"
  description = "Resource Group Name"
}

variable "cluster_name" {
  type        = string
  default     = "multi-cloud-aks-prod"
  description = "AKS Cluster Name"
}

variable "vnet_address_space" {
  type        = list(string)
  default     = ["10.0.0.0/16"]
  description = "VNet address CIDR block"
}

variable "aks_vm_size" {
  type        = string
  default     = "Standard_B2s"
  description = "AKS Node VM size (Standard_B2s or Standard_D2as_v4 for Spot)"
}
