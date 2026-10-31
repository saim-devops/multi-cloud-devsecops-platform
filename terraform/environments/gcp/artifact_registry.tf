resource "google_artifact_registry_repository" "core_repo" {
  location      = var.region
  repository_id = "core-registry"
  description   = "Multi-cloud Docker container registry"
  format        = "DOCKER"

  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"
    most_recent_versions {
      keep_count = 10
    }
  }
}
