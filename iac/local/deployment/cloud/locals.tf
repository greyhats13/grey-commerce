locals {
  region = var.region
  tags = {
    Unit       = var.unit
    Env        = var.env
    GithubRepo = var.github_repo
    GithubOrg  = var.github_owner
  }
  # AWS Secrets Manager Locals
  secrets_manager_standard = {
    Unit    = var.unit
    Env     = var.env
    Code    = "secretsmanager"
    Feature = "iac"
  }
  secrets_manager_naming_standard = "${local.secrets_manager_standard.Unit}/${local.secrets_manager_standard.Env}/${local.secrets_manager_standard.Code}/${local.secrets_manager_standard.Feature}"
  # ArgoCD Locals
  argocd_standard = {
    Unit    = var.unit
    Env     = var.env
    Code    = "addon"
    Feature = "argocd"
  }
  argocd_naming_standard = "${local.argocd_standard.Unit}-${local.argocd_standard.Env}-${local.argocd_standard.Code}-${local.argocd_standard.Feature}"
  # Github Locals
  github_action_variables = {
    DOCKER_USERNAME = var.docker_username
    GH_OWNER        = var.github_owner
    GH_REPO_NAME    = var.github_repo
  }
  github_action_secrets = {
    ARGOCD_SSH      = tls_private_key.argocd_ssh.private_key_pem
    DOCKER_PASSWORD = var.docker_password
  }
}
