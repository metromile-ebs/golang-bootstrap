# Example Generic Variables

variable "project" {
  type = string
  description = "common aws tags"
  default = "streamline-graph-manager"
}

variable "namespace" {
  type = string
  default = "common"
  description = "namespace for multi-tenant environment"
}

variable "environment" {
  type = string
  default = "dev"
  description = "type of environment <dev, stg, prod> for labeling"
}

variable "multi_tenant_default_tags" {
  type = map(string)
  description = "common aws tags"
  default = {
    "author"      = "terraform"
    "product"     = "streamline"
  }
}
