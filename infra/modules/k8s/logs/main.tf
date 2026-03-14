resource "helm_release" "prometheus" {
  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = var.chart_version
  namespace        = var.namespace
  create_namespace = true

  set {
    name  = "grafana.adminPassword"
    value = var.grafana_admin_password
  }

  set {
    name  = "grafana.service.type"
    value = "NodePort"
  }

  set {
    name  = "prometheus.prometheusSpec.retention"
    value = var.retention
  }

  dynamic "set" {
    for_each = var.values
    content {
      name  = split("=", set.value)[0]
      value = split("=", set.value)[1]
    }
  }
}
