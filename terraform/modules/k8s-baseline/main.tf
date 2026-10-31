terraform {
  required_version = ">= 1.5.0"
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.30"
    }
  }
}

resource "kubernetes_namespace" "core_service" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/part-of"          = "multi-cloud-devsecops"
      "pod-security.kubernetes.io/enforce" = "restricted"
      "pod-security.kubernetes.io/audit"   = "restricted"
      "pod-security.kubernetes.io/warn"    = "restricted"
    }
  }
}

resource "kubernetes_resource_quota" "core_service_quota" {
  metadata {
    name      = "${var.namespace}-quota"
    namespace = kubernetes_namespace.core_service.metadata[0].name
  }
  spec {
    hard = {
      "requests.cpu"    = "4"
      "requests.memory" = "8Gi"
      "limits.cpu"      = "8"
      "limits.memory"   = "16Gi"
      "pods"            = "20"
    }
  }
}

resource "kubernetes_limit_range" "core_service_limits" {
  metadata {
    name      = "${var.namespace}-limits"
    namespace = kubernetes_namespace.core_service.metadata[0].name
  }
  spec {
    limit {
      type = "Container"
      default = {
        cpu    = "250m"
        memory = "256Mi"
      }
      default_request = {
        cpu    = "50m"
        memory = "64Mi"
      }
    }
  }
}

resource "kubernetes_network_policy" "default_deny_all" {
  metadata {
    name      = "default-deny-ingress"
    namespace = kubernetes_namespace.core_service.metadata[0].name
  }
  spec {
    pod_selector {}
    policy_types = ["Ingress"]
  }
}

resource "kubernetes_network_policy" "allow_http_and_monitoring" {
  metadata {
    name      = "allow-http-and-monitoring"
    namespace = kubernetes_namespace.core_service.metadata[0].name
  }
  spec {
    pod_selector {
      match_labels = {
        "app.kubernetes.io/name" = "core-service"
      }
    }
    ingress {
      ports {
        port     = "8080"
        protocol = "TCP"
      }
    }
    policy_types = ["Ingress"]
  }
}
