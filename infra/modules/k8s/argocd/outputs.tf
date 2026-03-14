output "namespace" {
  description = "Kubernetes namespace where ArgoCD is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "ArgoCD Helm release name"
  value       = helm_release.argocd.name
}

output "app_of_apps_release_name" {
  description = "App-of-apps Helm release name"
  value       = helm_release.app_of_apps.name
}
