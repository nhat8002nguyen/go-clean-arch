# variables.tf
variable "aws_region" {
  description = "The AWS region to deploy to"
  type        = string
  default     = "ap-southeast-1"
}

variable "rds_instance_class" {
  description = "The instance class for RDS"
  type        = string
  default     = "db.t3.micro"
}

variable "db_name" {
  description = "The database name"
  type        = string
  default     = "ecommerce_app"
}

variable "db_username" {
  description = "The database username"
  type        = string
  default     = "nhatnguyen"
}

variable "db_password" {
  description = "The database password"
  type        = string
  default     = "your-db-password"
}

variable "vpc_id" {
  description = "The VPC ID to deploy ECS and RDS into"
  type        = string
}

variable "subnets" {
  description = "The list of subnets to deploy ECS and RDS into"
  type        = list(string)
}
