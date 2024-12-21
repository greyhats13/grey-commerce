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
    REDIS_PASSWORD = jsondecode(data.aws_secretsmanager_secret_version.iac.secret_string)["REDIS_PASSWORD"]
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
  autoscaling_enabled = true

  # Konfigurasi Autoscaling untuk Tabel Utama
  autoscaling_read = {
    scale_in_cooldown  = 50
    scale_out_cooldown = 40
    target_value       = 75
    max_capacity       = 20
    min_capacity       = 5
  }

  autoscaling_write = {
    scale_in_cooldown  = 50
    scale_out_cooldown = 40
    target_value       = 75
    max_capacity       = 20
    min_capacity       = 5
  }

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
    # Tambahkan atribut lain sesuai kebutuhan
  ]

  # Definisi Global Secondary Indexes (GSIs)
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

  # Konfigurasi Autoscaling untuk GSIs
  autoscaling_indexes = {
    EmailIndex = {
      read_min_capacity  = 5
      read_max_capacity  = 20
      write_min_capacity = 5
      write_max_capacity = 20
      target_value       = 75
    },
    CreatedAtIndex = {
      read_min_capacity  = 5
      read_max_capacity  = 20
      write_min_capacity = 5
      write_max_capacity = 20
      target_value       = 75
    },
    UpdatedAtIndex = {
      read_min_capacity  = 5
      read_max_capacity  = 20
      write_min_capacity = 5
      write_max_capacity = 20
      target_value       = 75
    }
  }

  # Tagging untuk pengelolaan
  tags = {
    Environment = "local"
    Service     = "user-service"
  }
}

# # Prepare Gitthub
# module "github_action_env" {
#   source                  = "../../../modules/github"
#   repo_name               = var.github_repo
#   owner                   = var.github_owner
#   svc_name                = local.svc_naming_standard
#   github_action_variables = local.github_action_variables
#   github_action_secrets   = local.github_action_secrets
# }
