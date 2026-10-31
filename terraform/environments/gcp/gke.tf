resource "google_service_account" "gke_node_sa" {
  account_id   = "${var.cluster_name}-node-sa"
  display_name = "GKE Node Service Account with Least Privilege"
}

resource "google_project_iam_member" "node_sa_roles" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
    "roles/artifactregistry.reader"
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.gke_node_sa.email}"
}

resource "google_container_cluster" "primary" {
  name     = var.cluster_name
  location = "${var.region}-a" # Single zone for cost reduction

  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  ip_allocation_policy {
    cluster_secondary_range_name  = "pod-ip-range"
    services_secondary_range_name = "service-ip-range"
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "172.16.0.0/28"
  }

  release_channel {
    channel = "REGULAR"
  }
}

# Cost-optimized Preemptible / Spot Node Pool
resource "google_container_node_pool" "spot_nodes" {
  name       = "${var.cluster_name}-spot-pool"
  location   = "${var.region}-a"
  cluster    = google_container_cluster.primary.name
  node_count = 2

  autoscaling {
    min_node_count = 1
    max_node_count = 3
  }

  node_config {
    preemptible  = true
    machine_type = "e2-standard-2"
    spot         = true

    service_account = google_service_account.gke_node_sa.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    labels = {
      environment = "production"
      node_type   = "spot"
    }

    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
}

# Workload Identity IAM binding for core-service
resource "google_service_account" "core_service_gsa" {
  account_id   = "core-service-sa"
  display_name = "Workload Identity SA for core-service microservice"
}

resource "google_service_account_iam_member" "workload_identity_user" {
  service_account_id = google_service_account.core_service_gsa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[core-service/core-service]"
}
