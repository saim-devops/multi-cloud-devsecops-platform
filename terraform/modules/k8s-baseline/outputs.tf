output "namespace_name" {
  value       = kubernetes_namespace.core_service.metadata[0].name
  description = "Created Kubernetes namespace name"
}
