variable "namespace" {
  description = "Kubernetes namespace for the monitoring stack"
  type        = string
  default     = "monitoring"
}

variable "chart_version" {
  description = "kube-prometheus-stack Helm chart version"
  type        = string
  default     = "56.6.2"
}

variable "grafana_admin_password" {
  description = "Grafana admin password"
  type        = string
  sensitive   = true
  default     = "admin"
}

variable "retention" {
  description = "Prometheus data retention period"
  type        = string
  default     = "7d"
}

variable "values" {
  description = "Additional Helm values as key=value strings"
  type        = list(string)
  default     = []
}
