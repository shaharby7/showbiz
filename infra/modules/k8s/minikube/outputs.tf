output "cluster_name" {
  description = "Name of the Minikube cluster"
  value       = minikube_cluster.this.cluster_name
}

output "client_certificate" {
  description = "Client certificate for cluster authentication"
  value       = minikube_cluster.this.client_certificate
  sensitive   = true
}

output "client_key" {
  description = "Client key for cluster authentication"
  value       = minikube_cluster.this.client_key
  sensitive   = true
}

output "cluster_ca_certificate" {
  description = "CA certificate for cluster authentication"
  value       = minikube_cluster.this.cluster_ca_certificate
  sensitive   = true
}

output "host" {
  description = "Kubernetes API server host"
  value       = minikube_cluster.this.host
}
