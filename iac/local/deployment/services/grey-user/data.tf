data aws_caller_identity current {}

data "aws_secretsmanager_secret_version" "iac" {
  secret_id = "grey/local/secretsmanager/iac"
}
