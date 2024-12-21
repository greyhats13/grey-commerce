resource "tls_private_key" "argocd_ssh" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "random_password" "redis" {
  length           = 12
  override_special = "!#$%&*"
  min_lower        = 3
  min_upper        = 3
  min_numeric      = 3
  min_special      = 0
}