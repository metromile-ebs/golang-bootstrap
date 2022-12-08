# Example Postgres RDS Terraform for Dev

# variable "rds_type" {
#   type = string
#   default = "db.t3.small"
#   description = "aws rds instance type"
# }

# variable "db_password" {
#   type = string
#   default = ""
#   description = "database password"
# }

# variable "is_multi_az" {
#   type = bool
#   default = false
#   description = "flag to indicate whether to make something multi availability zone or not"
# }

# resource "aws_db_instance" "streamline-deployments-rds-instance" {
#   identifier = "${var.project}-${var.namespace}-${var.environment}-database"
#   username   = replace("${title(var.project)}${title(var.namespace)}${title(var.environment)}${title("master")}", "-", "")
#   db_name    = replace("${title(var.project)}${title(var.namespace)}${title(var.environment)}${title("rds")}", "-", "")
#   password   = var.db_password

#   allocated_storage       = 10
#   engine                  = "postgres"
#   engine_version          = "13.7"
#   instance_class          = var.rds_type
#   parameter_group_name    = "default.postgres13"
#   multi_az                = var.is_multi_az
#   backup_retention_period = 7
#   backup_window           = "12:00-12:30"
#   maintenance_window      = "thu:06:33-thu:07:03"
#   publicly_accessible     = false
#   skip_final_snapshot     = true # Remove for all non dev environments

#   vpc_security_group_ids = var.eks_security_group_ids
#   db_subnet_group_name   = aws_db_subnet_group.deployments_rds_subnet.name

#   tags =tags = merge(
#     var.multi_tenant_default_tags,
#     {
#       environment   = var.environment,
#       namespace     = var.namespace,
#       project       = var.project
#     }
#   )
# }
