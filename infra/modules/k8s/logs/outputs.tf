output "namespace" {
  description = "Kubernetes namespace where the monitoring stack is deployed"
  value       = var.namespace
}

output "grafana_service" {
  description = "Grafana service name for port-forwarding"
  value       = "prometheus-grafana"
}

output "prometheus_service" {
  description = "Prometheus service name for port-forwarding"
  value       = "prometheus-kube-prometheus-prometheus"
}
