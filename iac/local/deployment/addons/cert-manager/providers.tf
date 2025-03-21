terraform {
  # backend "s3" {
  #   bucket = "grey-dev-s3-tfstate"
  #   key    = "grey/deployment/addons/grey-local-deployment-addons-cert-manager.tfstate"
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
  }
}

# Create AWS provider
provider "aws" {
  region = local.region
  assume_role {
    role_arn = "arn:aws:iam::124456474132:role/iac"
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
