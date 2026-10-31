variable "monitoring_namespace" {
  type        = string
  default     = "monitoring"
  description = "Namespace for Prometheus and Grafana stack"
}

variable "prometheus_chart_version" {
  type        = string
  default     = "61.3.1"
  description = "kube-prometheus-stack chart version"
}

variable "grafana_admin_password" {
  type        = string
  default     = "AdminDevSecOpsSecure2026!"
  sensitive   = true
  description = "Admin password for Grafana dashboard"
}
