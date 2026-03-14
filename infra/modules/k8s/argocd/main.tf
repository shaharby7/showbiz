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
