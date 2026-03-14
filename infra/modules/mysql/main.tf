resource "helm_release" "mysql" {
  name             = "mysql"
  repository       = "https://charts.bitnami.com/bitnami"
  chart            = "mysql"
  namespace        = var.namespace
  create_namespace = true

  set {
    name  = "auth.rootPassword"
    value = var.root_password
  }

  set {
    name  = "auth.database"
    value = var.database
  }

  set {
    name  = "auth.username"
    value = var.user
  }

  set {
    name  = "auth.password"
    value = var.password
  }

  set {
    name  = "primary.persistence.size"
    value = var.storage_size
  }
}
