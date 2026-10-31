resource "azurerm_user_assigned_identity" "aks_identity" {
  name                = "id-aks-${var.cluster_name}"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.cluster_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "aks-devsecops"
  kubernetes_version  = "1.29"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks_identity.id]
  }

  # Azure CNI Networking for Enterprise Multi-Cluster Routing
  network_profile {
    network_plugin    = "azure"
    network_policy    = "azure"
    load_balancer_sku = "standard"
  }

  default_node_pool {
    name                = "systempool"
    node_count          = 1
    vm_size             = var.aks_vm_size
    vnet_subnet_id      = azurerm_subnet.aks_nodes.id
    os_disk_type        = "Managed"
    os_disk_size_gb     = 30
    enable_auto_scaling = false
  }

  # Secrets Store CSI Driver Add-On (Validates AZ-400 Key Vault Integration)
  key_vault_secrets_provider {
    secret_rotation_enabled  = true
    secret_rotation_interval = "2m"
  }

  tags = {
    Environment = "production"
    ManagedBy   = "Terraform"
  }
}

# Cost-Saving Spot Node Pool for Application Workloads
resource "azurerm_kubernetes_cluster_node_pool" "spot_workloads" {
  name                  = "spotworkload"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.aks.id
  vm_size               = var.aks_vm_size
  vnet_subnet_id        = azurerm_subnet.aks_nodes.id
  priority              = "Spot"
  eviction_policy       = "Delete"
  spot_max_price        = -1 # Pay up to on-demand price

  enable_auto_scaling = true
  min_count           = 1
  max_count           = 3
  node_count          = 1

  node_labels = {
    "kubernetes.azure.com/scalesetpriority" = "spot"
    "workload"                              = "microservices"
  }

  node_taints = [
    "kubernetes.azure.com/scalesetpriority=spot:NoSchedule"
  ]
}
