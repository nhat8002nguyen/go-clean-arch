provider "aws" {
  region = var.aws_region
}

resource "aws_db_instance" "ecommerce_app_db" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "postgres"
  engine_version       = "16.3-R2"
  instance_class       = var.rds_instance_class
  identifier           = "ecommerce-app-db"
  username             = var.db_username
  password             = var.db_password
  parameter_group_name = "default.postgres16"
  skip_final_snapshot  = true

  vpc_security_group_ids = [aws_security_group.rds_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.ecommerce_app_db_subnet_group.name
}

resource "aws_db_subnet_group" "ecommerce_app_db_subnet_group" {
  name       = "ecommerce-app-db-subnet-group"
  subnet_ids = var.subnets
}

resource "aws_security_group" "rds_sg" {
  name   = "rds-security-group"
  vpc_id = var.vpc_id

  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Replace this with the actual SG or CIDR range for your ECS
  }
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
    subnets          = var.subnets
    security_groups  = [aws_security_group.ecs_sg.id]
    assign_public_ip = true
  }
}

resource "aws_security_group" "ecs_sg" {
  name   = "ecs-security-group"
  vpc_id = var.vpc_id

  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }

  ingress {
    from_port   = 9090
    to_port     = 9090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Update to actual access CIDR or Security Group
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
