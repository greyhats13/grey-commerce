module "iam_assumable_role_inline_policy" {
  source = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"

  trusted_role_arns = [
    "arn:aws:iam::000000000000:root",
  ]

  trusted_role_services = [
    "ec2.amazonaws.com"
  ]

  create_role = true

  role_name         = "external-secrets-role"
  role_requires_mfa = false
  inline_policy_statements = [
    {
      sid = "AllowECRPushPull"
      actions = [
        "secretsmanager:GetResourcePolicy",
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret",
        "secretsmanager:ListSecretVersionIds",
        "secretsmanager:CreateSecret",
        "secretsmanager:PutSecretValue",
        "secretsmanager:TagResource",
        "secretsmanager:DeleteSecret"
      ]
      effect    = "Allow"
      resources = ["*"]
    }
  ]
}

module "argocd_app" {
  source        = "../../../modules/helm"
  region        = var.region
  standard      = local.addon_standard
  override_name = local.addon_standard.Feature
  repository    = "https://argoproj.github.io/argo-helm"
  chart         = "argocd-apps"
  values        = ["${file("manifest/${local.addon_standard.Feature}.yaml")}"]
  namespace     = "argocd"
  dns_name      = "${local.addon_standard.Feature}.${var.unit}.blast.co.id"
  extra_vars = {
    argocd_namespace                       = "argocd"
    source_origin_repoURL                  = "https://kubernetes-sigs.github.io/external-secrets/charts"
    source_origin_chart                    = local.addon_standard.Feature
    source_origin_targetRevision           = "0.11.0"
    source_override_repoURL                = "git@github.com:${var.github_owner}/${var.github_repo}.git"
    source_override_targetRevision         = "main"
    source_override_path                   = "charts/local/addons/${local.addon_standard.Feature}/values.yaml"
    project                                = "default"
    destination_server                     = "https://kubernetes.default.svc"
    destination_namespace                  = local.addon_standard.Feature
    syncPolicy_automated_prune             = true
    syncPolicy_automated_selfHeal          = true
    syncPolicy_syncOptions_CreateNamespace = true
    cluster_secret_store_path              = "charts/local/addons/${local.addon_standard.Feature}/manifest/cluster-secret-store"
  }
}

resource "kubernetes_secret_v1" "secrets" {
  metadata {
    name      = "aws-creds"
    namespace = local.addon_standard.Feature
  }

  data = {
    access-key = "test"
    secret-key = "test"
  }
  depends_on = [module.argocd_app]
}
