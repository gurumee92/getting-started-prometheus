provider "aws" {
  region = "ap-northeast-2"
}

locals {
  cluster_name = "gurumee-test"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "gurumee-test"
  cidr = "10.194.0.0/16"

  azs             = ["ap-northeast-2a", "ap-northeast-2b", "ap-northeast-2c"]
  public_subnets  = ["10.194.0.0/24", "10.194.1.0/24", "10.194.2.0/24"]
  private_subnets = ["10.194.100.0/24", "10.194.101.0/24", "10.194.102.0/24"]

  enable_nat_gateway     = true
  one_nat_gateway_per_az = true

  enable_dns_hostnames = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"             = "1"
  }
}
