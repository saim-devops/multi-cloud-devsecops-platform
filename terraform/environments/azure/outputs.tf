output "aks_cluster_name" {
  value       = azurerm_kubernetes_cluster.aks.name
  description = "AKS Cluster name"
}

output "acr_login_server" {
  value       = azurerm_container_registry.acr.login_server
  description = "Azure Container Registry login server"
}

output "key_vault_uri" {
  value       = azurerm_key_vault.kv.vault_uri
  description = "Azure Key Vault URI"
}

output "key_vault_csi_client_id" {
  value       = azurerm_kubernetes_cluster.aks.key_vault_secrets_provider[0].secret_identity[0].client_id
  description = "Client ID of the AKS Key Vault Secrets Provider Managed Identity"
}

output "kubeconfig_command" {
  value       = "az aks get-credentials --resource-group ${azurerm_resource_group.rg.name} --name ${azurerm_kubernetes_cluster.aks.name} --overwrite-existing"
  description = "CLI command to retrieve AKS kubeconfig"
}
