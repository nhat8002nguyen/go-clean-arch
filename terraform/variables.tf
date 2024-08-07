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

variable "context_timeout" {
  description = "The application context timeout"
  type        = number
}

variable "server_address" {
  description = "The application server address port"
  type        = string
}

variable "is_debug" {
  description = "The debug mode"
  type        = bool
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
}

variable "db_driver" {
  description = "The database driver"
  type        = string
}

variable "local_ips" {
  description = "The local ip addresses that can access the bastion host."
  type        = list(string)
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "The secret access key to AWS"
  type        = string
}

variable "AWS_ACCESS_KEY_ID" {
  description = "the access key id to AWS"
  type        = string
}
