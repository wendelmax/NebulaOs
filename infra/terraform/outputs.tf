output "environment" {
  value = var.environment
}

output "project_id" {
  value = "nebula-${var.environment}"
}
