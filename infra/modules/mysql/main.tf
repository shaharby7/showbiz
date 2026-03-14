terraform {
  required_version = ">= 1.0"
}

variable "namespace" {
  description = "Kubernetes namespace for MySQL"
  type        = string
  default     = "showbiz"
}

variable "root_password" {
  description = "MySQL root password"
  type        = string
  sensitive   = true
  default     = "rootpassword"
}

variable "database" {
  description = "Database name to create"
  type        = string
  default     = "showbiz"
}

variable "user" {
  description = "MySQL user"
  type        = string
  default     = "showbiz"
}

variable "password" {
  description = "MySQL user password"
  type        = string
  sensitive   = true
  default     = "showbiz_dev"
}

variable "storage_size" {
  description = "PVC storage size"
  type        = string
  default     = "5Gi"
}

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

output "host" {
  value = "mysql.${var.namespace}.svc.cluster.local"
}

output "port" {
  value = 3306
}

output "database" {
  value = var.database
}
