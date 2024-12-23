# Secrets Manager
module "secrets_iac" {
  source  = "terraform-aws-modules/secrets-manager/aws"
  version = "~> 1.3.1"

  # Secret
  name                    = local.secrets_manager_naming_standard
  description             = "Secrets for ${local.secrets_manager_naming_standard}"
  recovery_window_in_days = 0
  # Policy
  create_policy       = true
  block_public_policy = true
  policy_statements = {
    admin = {
      sid = "IacSecretAdmin"
      principals = [
        {
          type        = "AWS"
          identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
        }
      ]
      actions   = ["secretsmanager:GetSecretValue"]
      resources = ["*"]
    }
  }

  # Version
  ignore_secret_changes = true
  secret_string = jsonencode({
    argocd_ssh_base64 = base64encode(tls_private_key.argocd_ssh.private_key_pem)
    REDIS_PASSWORD    = random_password.redis.result
  })

  tags = merge(local.tags, local.secrets_manager_standard)
}

# ArgoCD
# ArgoCD is a declarative, GitOps continuous delivery tool for Kubernetes.
module "argocd" {
  source = "../../modules/helm"

  region           = var.region
  standard         = local.argocd_standard
  override_name    = local.argocd_standard.Feature
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  values           = ["${file("manifest/${local.argocd_standard.Feature}.yaml")}"]
  namespace        = local.argocd_standard.Feature
  create_namespace = true
  dns_name         = var.local_dns
  extra_vars = {
    github_orgs      = var.github_orgs
    github_client_id = var.github_oauth_client_id
    ARGOCD_VERSION   = var.argocd_version
    AVP_VERSION      = var.argocd_vault_plugin_version
    server_insecure  = false

    # ref https://github.com/argoproj/argo-helm/tree/main/charts/argo-cd
    # ingress
    ingress_enabled    = true
    ingress_class_name = "nginx"
  }
}

# Setup repository for argocd and atlantis
module "github" {
  source                     = "../../modules/github"
  repo_name                  = var.github_repo
  owner                      = var.github_owner
  create_deploy_key          = true
  add_repo_ssh_key_to_argocd = true
  public_key                 = tls_private_key.argocd_ssh.public_key_openssh
  ssh_key                    = tls_private_key.argocd_ssh.private_key_pem
  is_deploy_key_read_only    = false
  argocd_namespace           = "argocd"
  github_action_variables    = local.github_action_variables
  github_action_secrets      = local.github_action_secrets
  depends_on = [
    module.argocd,
  ]
}
