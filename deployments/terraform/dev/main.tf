
# Example Terraform Main File for MT Dev
terraform {
  backend "s3" {
    bucket = "streamline-deployments"
    key = "${var.project}-service/dev"
    region = "us-west-2"
  }
  required_version = "1.0.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-west-2"
}

module "streamline-graph-manager-module" {
  source = "../modules"

  namespace   = "common"
  environment = "dev"

  nginx_ingress_cname = "a8eaeff97dd1d463caaa70ffe404a5a9-278258945.us-west-2.elb.amazonaws.com"
}
