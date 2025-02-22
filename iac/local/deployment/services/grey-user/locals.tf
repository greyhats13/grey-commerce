locals {
  region = var.region
  tags = {
    GithubRepo = var.github_repo
    GithubOrg  = var.github_owner
  }
  svc_standard = {
    Unit    = var.unit
    Env     = var.env
    Code    = "svc"
    Feature = "user"
  }
  svc_naming_standard = "${local.svc_standard.Unit}-${local.svc_standard.Code}-${local.svc_standard.Feature}"
  svc_naming_full     = "${local.svc_standard.Unit}-${local.svc_standard.Env}-${local.svc_standard.Code}-${local.svc_standard.Feature}"
  svc_name            = "${local.svc_standard.Unit}-${local.svc_standard.Feature}"
  svc_secret_standard = "${local.svc_standard.Unit}/${local.svc_standard.Code}/${local.svc_standard.Feature}"
  ## Environment variables that will be stored in Github repo environment for Github Actions
  github_action_variables_env = {
    svc_name            = local.svc_name
    svc_naming_standard = local.svc_naming_standard
    svc_naming_full     = local.svc_naming_full
    gitops_path_local   = "local"
    gitops_path_dev     = "incubator"
    gitops_path_stg     = "test"
    gitops_path_prod    = "stable"
  }
  ## Environment secrets that will be stored in Github repo environment for Github Actions
  github_action_secrets_env = {
    argocd_ssh = base64decode(jsondecode(data.aws_secretsmanager_secret_version.iac.secret_string)["argocd_ssh_base64"])
  }
}
