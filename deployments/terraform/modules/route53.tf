# Example Route53 Ingress Terraform
# Requires existing nginx ingress to exist

# variable "nginx_ingress_cname" {
#   type = string
#   default = ""
#   description = "Nginx Ingress DNS Record"
# }

# data "aws_route53_zone" "primary" {
#   name = "metromileai.com"
#   private_zone = false
# }

# resource "aws_route53_record" "deployments_route53_record" {
#   zone_id = data.aws_route53_zone.primary.zone_id
#   name    = "${var.project}-${var.namespace}-${var.environment}.metromileai.com"
#   type    = "CNAME"
#   ttl     = "30"
#   records = [
#     var.nginx_ingress_cname
#   ]
# }
