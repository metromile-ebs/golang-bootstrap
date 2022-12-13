variable "eks_security_group_ids" {
  type = list(string)
  default = []
  description = "list of eks security group ids"
}
variable "private_subnet_ids" {
  type = list(string)
  default = []
  description = "private subnet ids used for RDS"
}

variable "private_subnet_availability_zones" {
  type = list(string)
  default = []
  description = "list of availablity zones for private subnets"
}

resource "aws_db_subnet_group" "graphs_docdb_subnet" {
  name        = "${var.project}-${var.namespace}-${var.environment}-subnet"
  subnet_ids  = var.private_subnet_ids
  tags        = merge(
    var.multi_tenant_default_tags,
    {
      environment   = var.environment,
      namespace     = var.namespace,
      project       = var.project
    }
  )
}