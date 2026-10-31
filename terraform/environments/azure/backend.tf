terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.100"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.48"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.30"
    }
  }

  # Production remote state backend with Azure Blob Storage:
  # backend "azurerm" {
  #   resource_group_name  = "rg-devsecops-tfstate"
  #   storage_account_name = "stdevsecopstfstate"
  #   container_name       = "tfstate"
  #   key                  = "azure.terraform.tfstate"
  # }
}

provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}
