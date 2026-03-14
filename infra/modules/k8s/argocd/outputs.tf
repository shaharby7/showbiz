output "namespace" {
  description = "Kubernetes namespace where ArgoCD is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = helm_release.argocd.name
}
