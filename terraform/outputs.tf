# outputs.tf
output "ecs_cluster_id" {
  value = aws_ecs_cluster.ecommerce_app_cluster.id
}

output "rds_endpoint" {
  value = aws_db_instance.ecommerce_app_db.endpoint
}

output "rds_host" {
  value = aws_db_instance.ecommerce_app_db.address
}

output "rds_port" {
  value = aws_db_instance.ecommerce_app_db.port
}

output "rds_db_name" {
  value = var.db_name
}
