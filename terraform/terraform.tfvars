# terraform.tfvars
aws_region         = "ap-southeast-1"
rds_instance_class = "db.t3.micro"
is_debug           = false
context_timeout    = 3
server_address     = ":9090"
db_name            = "ecommerce_app"
db_username        = "nhatnguyen"
db_password        = "samplepassword"
db_driver          = "postgres"
local_ips          = ["42.115.164.145/32"]
