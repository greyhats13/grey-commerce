terraform {
  backend "s3" {
    bucket = "grey-dev-s3-tfstate"
    key    = "grey/deployment/svc/grey-dev-deployment-svc-profile.tfstate"
    region = "ap-southeast-1"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.81.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.34.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.16.1"
    }
    mysql = {
      source  = "petoju/mysql"
      version = "3.0.67"
    }
    github = {
      source  = "integrations/github"
      version = "~> 6.4.0"
    }
  }
}

# Create AWS provider
provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = "ap-southeast-1"
  s3_use_path_style           = false
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = "https://localstack.lokal.blast.co.id6"
    apigatewayv2   = "https://localstack.lokal.blast.co.id6"
    cloudformation = "https://localstack.lokal.blast.co.id6"
    cloudwatch     = "https://localstack.lokal.blast.co.id6"
    dynamodb       = "https://localstack.lokal.blast.co.id6"
    ec2            = "https://localstack.lokal.blast.co.id6"
    es             = "https://localstack.lokal.blast.co.id6"
    elasticache    = "https://localstack.lokal.blast.co.id6"
    firehose       = "https://localstack.lokal.blast.co.id6"
    iam            = "https://localstack.lokal.blast.co.id6"
    kinesis        = "https://localstack.lokal.blast.co.id6"
    lambda         = "https://localstack.lokal.blast.co.id6"
    rds            = "https://localstack.lokal.blast.co.id6"
    redshift       = "https://localstack.lokal.blast.co.id6"
    route53        = "https://localstack.lokal.blast.co.id6"
    s3             = "http://s3.localhost.localstack.cloud:4566"
    secretsmanager = "https://localstack.lokal.blast.co.id6"
    ses            = "https://localstack.lokal.blast.co.id6"
    sns            = "https://localstack.lokal.blast.co.id6"
    sqs            = "https://localstack.lokal.blast.co.id6"
    ssm            = "https://localstack.lokal.blast.co.id6"
    stepfunctions  = "https://localstack.lokal.blast.co.id6"
    sts            = "https://localstack.lokal.blast.co.id6"
  }
}

# Create Kubernetes provider
provider "kubernetes" {
  host                   = data.terraform_remote_state.cloud.outputs.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(data.terraform_remote_state.cloud.outputs.eks_cluster_certificate_authority_data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

# Create Helm provider
provider "helm" {
  kubernetes {
    host                   = data.terraform_remote_state.cloud.outputs.eks_cluster_endpoint
    cluster_ca_certificate = base64decode(data.terraform_remote_state.cloud.outputs.eks_cluster_certificate_authority_data)
    token                  = data.aws_eks_cluster_auth.cluster.token
  }
}

provider "mysql" {
  endpoint = "${data.terraform_remote_state.cloud.outputs.aurora_cluster_endpoint}:${data.terraform_remote_state.cloud.outputs.aurora_cluster_port}"
  username = data.terraform_remote_state.cloud.outputs.aurora_cluster_username
  password = jsondecode(data.aws_secretsmanager_secret_version.aurora_password.secret_string)["password"]
}
