variable "db_password" {
    type = string
    description = "docdb password"
}

resource "aws_docdb_cluster" "docdb" {
  cluster_identifier      = "${var.project}-${var.namespace}-${var.environment}-cluster"
  engine                  = "docdb"
  master_username         = "${var.project}-${var.namespace}-${var.environment}-master"
  master_password         = var.db_password

  apply_immediately       = true
  backup_retention_period = 5
  preferred_backup_window = "00:00-01:00"
  skip_final_snapshot     = false

  availability_zones      = var.private_subnet_availability_zones
  db_subnet_group_name    = aws_db_subnet_group.graphs_docdb_subnet.name
  vpc_security_group_ids  = var.eks_security_group_ids

  tags = merge(
    var.multi_tenant_default_tags,
    {
      environment   = var.environment,
      namespace     = var.namespace,
      project       = var.project
    }
  )
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  count              = 1
  identifier         = "${var.projectName}-${var.namespace}-${var.environment}-instance"
  cluster_identifier = aws_docdb_cluster.docdb.id
  instance_class     = "db.t4g.medium"
}