variable "namespace" {
  description = "Kubernetes namespace for ArgoCD"
  type        = string
  default     = "argocd"
}

variable "chart_version" {
  description = "ArgoCD Helm chart version"
  type        = string
  default     = "5.51.6"
}

variable "argocd_values" {
  description = "Additional ArgoCD Helm values as key=value strings"
  type        = list(string)
  default     = []
}

variable "environment" {
  description = "Environment name (local, staging, production) — passed to the app-of-apps chart"
  type        = string
}

variable "repo_url" {
  description = "Git repository URL for the Showbiz monorepo"
  type        = string
}

variable "target_revision" {
  description = "Git revision for ArgoCD to track (branch, tag, or SHA)"
  type        = string
  default     = "HEAD"
}

variable "app_of_apps_chart_path" {
  description = "Local path to the app-of-apps Helm chart"
  type        = string
}
