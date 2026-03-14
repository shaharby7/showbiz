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
