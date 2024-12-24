terraform {
  # backend "s3" {
  #   bucket = "grey-dev-s3-tfstate"
  #   key    = "grey/deployment/svc/grey-dev-deployment-svc-user.tfstate"
  #   region = "ap-southeast-1"
  # }
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
  region                      = var.region
  s3_use_path_style           = false
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = "https://localstack.lokal.blast.co.id"
    apigatewayv2   = "https://localstack.lokal.blast.co.id"
    cloudformation = "https://localstack.lokal.blast.co.id"
    cloudwatch     = "https://localstack.lokal.blast.co.id"
    dynamodb       = "https://localstack.lokal.blast.co.id"
    ec2            = "https://localstack.lokal.blast.co.id"
    es             = "https://localstack.lokal.blast.co.id"
    elasticache    = "https://localstack.lokal.blast.co.id"
    firehose       = "https://localstack.lokal.blast.co.id"
    iam            = "https://localstack.lokal.blast.co.id"
    kinesis        = "https://localstack.lokal.blast.co.id"
    lambda         = "https://localstack.lokal.blast.co.id"
    rds            = "https://localstack.lokal.blast.co.id"
    redshift       = "https://localstack.lokal.blast.co.id"
    route53        = "https://localstack.lokal.blast.co.id"
    s3             = "http://s3.localhost.localstack.cloud:4566"
    secretsmanager = "https://localstack.lokal.blast.co.id"
    ses            = "https://localstack.lokal.blast.co.id"
    sns            = "https://localstack.lokal.blast.co.id"
    sqs            = "https://localstack.lokal.blast.co.id"
    ssm            = "https://localstack.lokal.blast.co.id"
    stepfunctions  = "https://localstack.lokal.blast.co.id"
    sts            = "https://localstack.lokal.blast.co.id"
  }
}

# Create Kubernetes provider
provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "docker-desktop"
}

# Create Helm provider
provider "helm" {
  kubernetes {
    config_path    = "~/.kube/config"
    config_context = "docker-desktop"
  }
}

