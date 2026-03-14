terraform {
  required_version = ">= 1.0"
}

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

variable "values" {
  description = "Additional Helm values"
  type        = list(string)
  default     = []
}

resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  version          = var.chart_version
  namespace        = var.namespace
  create_namespace = true

  set {
    name  = "server.service.type"
    value = "NodePort"
  }

  dynamic "set" {
    for_each = var.values
    content {
      name  = split("=", set.value)[0]
      value = split("=", set.value)[1]
    }
  }
}

output "namespace" {
  value = var.namespace
}
