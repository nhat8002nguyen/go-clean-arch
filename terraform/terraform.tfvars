# terraform.tfvars
aws_region         = "ap-southeast-1"
rds_instance_class = "db.t3.micro"
db_name            = "ecommerce_app"
db_username        = "nhatnguyen"
db_password        = "123456@RDS"
vpc_id             = "vpc-0a350f5c25581d781"
subnets            = ["subnet-093d57c58fd57e105", "subnet-01a5c992febe62f64", "subnet-0cae8ee564af5726a"]
