output "cluster_name" {
  value       = google_container_cluster.primary.name
  description = "GKE Cluster name"
}

output "cluster_endpoint" {
  value       = google_container_cluster.primary.endpoint
  description = "GKE Cluster API Endpoint"
}

output "artifact_registry_repo" {
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.core_repo.name}"
  description = "Google Artifact Registry Docker URI"
}

output "workload_identity_sa" {
  value       = google_service_account.core_service_gsa.email
  description = "Google Service Account bound to K8s via Workload Identity"
}

output "kubeconfig_command" {
  value       = "gcloud container clusters get-credentials ${google_container_cluster.primary.name} --zone ${var.region}-a --project ${var.project_id}"
  description = "CLI command to retrieve GKE credentials"
}
