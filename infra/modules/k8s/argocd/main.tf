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
    for_each = var.argocd_values
    content {
      name  = split("=", set.value)[0]
      value = split("=", set.value)[1]
    }
  }
}

resource "helm_release" "app_of_apps" {
  name      = "showbiz-app-of-apps"
  chart     = var.app_of_apps_chart_path
  namespace = var.namespace

  depends_on = [helm_release.argocd]

  set {
    name  = "environment"
    value = var.environment
  }

  set {
    name  = "repoURL"
    value = var.repo_url
  }

  set {
    name  = "targetRevision"
    value = var.target_revision
  }
}
