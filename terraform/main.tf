provider "aws" {
  region = var.aws_region
}

# VPC and Subnets
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "main-vpc"
  }
  enable_dns_hostnames = true
  enable_dns_support   = true
}

# Private Subnets for RDS
resource "aws_subnet" "private_subnet_rds_1" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.2.0/24"     # Adjust CIDR block as needed
  availability_zone = "ap-southeast-1a" # Replace with your AZ
  tags = {
    Name = "rds-private-subnet-1"
  }
}

# Private Subnets for RDS
resource "aws_subnet" "private_subnet_rds_2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.3.0/24"     # Adjust CIDR block as needed
  availability_zone = "ap-southeast-1b" # Replace with another AZ
  tags = {
    Name = "rds-private-subnet-2"
  }
}

# Private Subnets for RDS
resource "aws_subnet" "private_subnet_rds_3" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.4.0/24"     # Adjust CIDR block as needed
  availability_zone = "ap-southeast-1c" # Replace with another AZ
  tags = {
    Name = "rds-private-subnet-3"
  }
}

# Route Table for Private Subnets (if not already created)
resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "rds-private-route-table"
  }
}

# Route Table Associations for Private Subnets (if not already created)
resource "aws_route_table_association" "private_subnet_association_1" {
  subnet_id      = aws_subnet.private_subnet_rds_1.id
  route_table_id = aws_route_table.private_route_table.id
}

resource "aws_route_table_association" "private_subnet_association_2" {
  subnet_id      = aws_subnet.private_subnet_rds_2.id
  route_table_id = aws_route_table.private_route_table.id
}

resource "aws_route_table_association" "private_subnet_association_3" {
  subnet_id      = aws_subnet.private_subnet_rds_3.id
  route_table_id = aws_route_table.private_route_table.id
}

resource "aws_db_subnet_group" "ecommerce_app_db_subnet_group" {
  name       = "ecommerce-app-db-subnet-group"
  subnet_ids = [aws_subnet.private_subnet_rds_1.id, aws_subnet.private_subnet_rds_2.id, aws_subnet.private_subnet_rds_3.id]
}

resource "aws_security_group" "rds_sg" {
  name   = "rds-security-group"
  vpc_id = aws_vpc.main.id

  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_sg.id]
  }
}

resource "aws_db_instance" "ecommerce_app_db" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "postgres"
  engine_version       = "16.3"
  instance_class       = var.rds_instance_class
  identifier           = "ecommerce-app-db"
  username             = var.db_username
  password             = var.db_password
  parameter_group_name = "default.postgres16"
  skip_final_snapshot  = true
  publicly_accessible  = true

  vpc_security_group_ids = [aws_security_group.rds_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.ecommerce_app_db_subnet_group.name
}

resource "aws_ecs_cluster" "ecommerce_app_cluster" {
  name = "ecommerce-app-cluster"
}

resource "aws_ecs_task_definition" "ecommerce_app_task" {
  family                   = "ecommerce-app-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([{
    name      = "ecommerce-app-container"
    image     = "700876988155.dkr.ecr.ap-southeast-1.amazonaws.com/ecommerce-go-app:latest"
    cpu       = 256
    memory    = 512
    essential = true
    portMappings = [{
      containerPort = 9090
      hostPort      = 9090
    }]
    environment = [
      {
        name  = "POSTGRES_DB"
        value = var.db_name
      },
      {
        name  = "POSTGRES_USER"
        value = var.db_username
      },
      {
        name  = "POSTGRES_PASSWORD"
        value = var.db_password
      },
      {
        name  = "POSTGRES_HOST"
        value = aws_db_instance.ecommerce_app_db.address
      }
    ]
  }])
}

resource "aws_ecs_service" "ecommerce_app_service" {
  name            = "ecommerce-app-service"
  cluster         = aws_ecs_cluster.ecommerce_app_cluster.id
  task_definition = aws_ecs_task_definition.ecommerce_app_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets = [
      aws_subnet.private_subnet_rds_1.id,
      aws_subnet.private_subnet_rds_2.id,
      aws_subnet.private_subnet_rds_3.id
    ]
    security_groups  = [aws_security_group.ecs_sg.id, aws_security_group.bastion.id]
    assign_public_ip = false
  }
}

resource "aws_security_group" "ecs_sg" {
  name   = "ecs-security-group"
  vpc_id = aws_vpc.main.id

  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecs_task_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
  ]
}

# This subnet is used for bastion host to connect to out private ECS, RDS.
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "ap-southeast-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "bastion-public-subnet"
  }
}

# Internet Gateway for Public Subnet
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "bastion-internet-gateway"
  }
}

# Route Table for Public Subnet
resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }

  tags = {
    Name = "bastion-public-route-table"
  }
}

# Route Table Association for Public Subnet
resource "aws_route_table_association" "public_subnet_association" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public_route_table.id
}

# Security Groups
resource "aws_security_group" "bastion" {
  name        = "bastion-security-group"
  description = "Security group for bastion host"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.local_ips # Replace with your IP
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Bastion Host
resource "aws_instance" "bastion" {
  ami                    = "ami-060e277c0d4cce553" # Ubuntu Server 24.04 LTS in ap-southeast-1
  instance_type          = "t2.micro"
  key_name               = "ecommerce_app" # Replace with your key pair
  subnet_id              = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.bastion.id]

  tags = {
    Name = "bastion-host"
  }
}
