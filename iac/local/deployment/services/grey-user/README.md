<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| aws | 5.81.0 |
| github | ~> 6.4.0 |
| helm | 2.16.1 |
| kubernetes | 2.34.0 |

## Providers

| Name | Version |
|------|---------|
| aws | 5.81.0 |
| random | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| argocd\_app | ../../../modules/helm | n/a |
| dynamodb\_table | terraform-aws-modules/dynamodb-table/aws | n/a |
| github\_action\_env | ../../../modules/github | n/a |
| secrets\_manager | terraform-aws-modules/secrets-manager/aws | ~> 1.3.1 |

## Resources

| Name | Type |
|------|------|
| [random_password.password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/5.81.0/docs/data-sources/caller_identity) | data source |
| [aws_secretsmanager_secret_version.iac](https://registry.terraform.io/providers/hashicorp/aws/5.81.0/docs/data-sources/secretsmanager_secret_version) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| env | Stage environment | `string` | n/a | yes |
| github\_owner | Github repository owner | `string` | n/a | yes |
| github\_repo | Github repository name | `string` | n/a | yes |
| region | AWS region | `string` | n/a | yes |
| unit | Business unit code | `string` | n/a | yes |
<!-- END_TF_DOCS -->