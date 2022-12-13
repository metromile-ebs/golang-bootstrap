
# Example Terraform Main File for MT Dev
# terraform {
#   backend "s3" {
#     bucket = "streamline-deployments"
#     key = "<INSERT SERVICE NAME>-service/dev"
#     region = "us-west-2"
#   }
#   required_version = "1.0.7"
#   required_providers {
#     aws = {
#       source  = "hashicorp/aws"
#       version = "~> 4.0"
#     }
#   }
# }

# provider "aws" {
#   region = "us-west-2"
# }

# module "streamline-deployments-module" {
#   source = "../modules"

#   namespace   = "common"
#   environment = "dev"

#   db_password = "abc123"
#   nginx_ingress_cname = "a8eaeff97dd1d463caaa70ffe404a5a9-278258945.us-west-2.elb.amazonaws.com"

#   private_subnet_ids = [
#     "subnet-0f797924f9dc89e77",
#     "subnet-02fc4bc08a8ad782d",
#     "subnet-029153d743a7775d5"
#   ]
#   private_subnet_availability_zones = [
#     "us-west-2b",
#     "us-west-2c",
#     "us-west-2d"
#   ]
#   eks_security_group_ids = [
#     "sg-04f317ed8a91c2361"
#   ]
# }
