output "endpoint" {
  description = "RDS instance endpoint (host:port)"
  value       = aws_db_instance.this.endpoint
}

output "host" {
  description = "RDS instance hostname"
  value       = aws_db_instance.this.address
}

output "port" {
  description = "RDS instance port"
  value       = aws_db_instance.this.port
}

output "database" {
  description = "Name of the created database"
  value       = aws_db_instance.this.db_name
}

output "arn" {
  description = "ARN of the RDS instance"
  value       = aws_db_instance.this.arn
}
