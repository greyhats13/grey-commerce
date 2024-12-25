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
  ignore_secret_changes = false
  secret_string = jsonencode({
    REDIS_PASSWORD = jsondecode(data.aws_secretsmanager_secret_version.iac.secret_string)["REDIS_PASSWORD"]
    AWS_ACCESS_KEY_ID = "test"
    AWS_SECRET_ACCESS_KEY = "test"
  })

  tags = merge(local.tags, local.svc_standard)
}

# Create DynamoDB Table
module "dynamodb_table" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  name                = local.svc_naming_standard
  hash_key            = "UUID"
  billing_mode        = "PROVISIONED"
  read_capacity       = 5
  write_capacity      = 5
  autoscaling_enabled = false

  # # Configure autoscaling for main table
  # autoscaling_read = {
  #   scale_in_cooldown  = 50
  #   scale_out_cooldown = 40
  #   target_value       = 75
  #   max_capacity       = 20
  #   min_capacity       = 5
  # }

  # autoscaling_write = {
  #   scale_in_cooldown  = 50
  #   scale_out_cooldown = 40
  #   target_value       = 75
  #   max_capacity       = 20
  #   min_capacity       = 5
  # }

  # Definisi Atribut Tabel Utama dan GSIs
  attributes = [
    {
      name = "UUID"
      type = "S"
    },
    {
      name = "Email"
      type = "S"
    },
    {
      name = "CreatedAt"
      type = "S"
    },
    {
      name = "UpdatedAt"
      type = "S"
    },

  ]

  # Define Global Secondary Indexes (GSIs)
  global_secondary_indexes = [
    {
      name            = "EmailIndex"
      hash_key        = "Email"
      projection_type = "ALL"
      write_capacity  = 5
      read_capacity   = 5
    },
    {
      name            = "CreatedAtIndex"
      hash_key        = "CreatedAt"
      sort_key        = "Email"
      projection_type = "ALL"
      write_capacity  = 5
      read_capacity   = 5
    },
    {
      name            = "UpdatedAtIndex"
      hash_key        = "UpdatedAt"
      sort_key        = "Email"
      projection_type = "ALL"
      write_capacity  = 5
      read_capacity   = 5
    }
  ]

  # Configure autoscaling for GSIs
  # autoscaling_indexes = {
  #   EmailIndex = {
  #     read_min_capacity  = 5
  #     read_max_capacity  = 20
  #     write_min_capacity = 5
  #     write_max_capacity = 20
  #     target_value       = 75
  #   },
  #   CreatedAtIndex = {
  #     read_min_capacity  = 5
  #     read_max_capacity  = 20
  #     write_min_capacity = 5
  #     write_max_capacity = 20
  #     target_value       = 75
  #   },
  #   UpdatedAtIndex = {
  #     read_min_capacity  = 5
  #     read_max_capacity  = 20
  #     write_min_capacity = 5
  #     write_max_capacity = 20
  #     target_value       = 75
  #   }
  # }

  # Tagging untuk pengelolaan
  tags = {
    Environment = "local"
    Service     = "grey-svc-user"
  }
}

# Prepare Gitthub
module "github_action_env" {
  source                      = "../../../modules/github"
  repo_name                   = var.github_repo
  owner                       = var.github_owner
  svc_name                    = local.svc_naming_standard
  github_action_variables_env = local.github_action_variables_env
  # github_action_secrets_env   = local.github_action_secrets_env
}

## Create ArgoCD App
module "argocd_app" {
  source     = "../../../modules/helm"
  region     = var.region
  standard   = local.svc_standard
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argocd-apps"
  values     = ["${file("manifest/${local.svc_standard.Feature}.yaml")}"]
  namespace  = "argocd"
  dns_name   = "${local.svc_standard.Feature}.${var.unit}.blast.co.id"
  extra_vars = {
    argocd_namespace      = "argocd"
    source_repoURL        = "git@github.com:${var.github_owner}/${var.github_repo}.git"
    source_targetRevision = "local"
    source_path = var.env == "dev" ? "charts/incubator/app/${local.svc_name}" : (
      var.env == "stg" ? "charts/test/app/${local.svc_name}" : (
        var.env == "prod" ? "charts/stable/app/${local.svc_name}" : "charts/local/app/${local.svc_name}"
      )
    )
    project                                = "default"
    destination_server                     = "https://kubernetes.default.svc"
    destination_namespace                  = var.env
    avp_type                               = "awssecretsmanager"
    region                                 = var.region
    syncPolicy_automated_prune             = true
    syncPolicy_automated_selfHeal          = true
    syncPolicy_syncOptions_CreateNamespace = true
  }
}
