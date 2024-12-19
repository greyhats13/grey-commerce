# Secrets Manager
## Create Secrets Manager
module "secrets_manager" {
  source  = "terraform-aws-modules/secrets-manager/aws"
  version = "~> 1.3.1"

  # Secret
  name                    = local.svc_secret_standard
  description             = "Secrets for ${local.svc_secret_standard}"
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
    db_host     = data.terraform_remote_state.cloud.outputs.aurora_cluster_endpoint
    db_port     = tostring(data.terraform_remote_state.cloud.outputs.aurora_cluster_port)
    db_name     = mysql_database.db.name
    db_user     = mysql_user.db.user
    db_password = random_password.password.result
  })

  tags = merge(local.tags, local.svc_standard)
}

# Prepare GIthub
module "github_action_env" {
  source                  = "../../../modules/github"
  repo_name               = var.github_repo
  owner                   = var.github_owner
  svc_name                = local.svc_naming_standard
  github_action_variables = local.github_action_variables
  # github_action_secrets   = local.github_action_secrets
}