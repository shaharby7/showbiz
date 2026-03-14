output "host" {
  description = "MySQL service hostname within the cluster"
  value       = "mysql.${var.namespace}.svc.cluster.local"
}

output "port" {
  description = "MySQL service port"
  value       = 3306
}

output "database" {
  description = "Name of the created database"
  value       = var.database
}
