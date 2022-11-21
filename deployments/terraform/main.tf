terraform {
  backend "s3" {
    bucket = "golang-bootstrap"
    key = "golang-bootstrap-service"
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
  region = "us-east-1"
}