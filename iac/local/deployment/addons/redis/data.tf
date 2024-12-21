data "aws_secretsmanager_secret_version" "redis" {
  secret_id = "grey/local/secretsmanager/iac"
}
