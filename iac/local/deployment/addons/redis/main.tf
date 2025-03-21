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
    source_origin_repoURL                  = "https://charts.bitnami.com/bitnami"
    source_origin_chart                    = "${local.addon_standard.Feature}"
    source_origin_targetRevision           = "17.16.0"
    source_override_repoURL                = "git@github.com:${var.github_owner}/${var.github_repo}.git"
    source_override_targetRevision         = "local"
    source_override_path                   = "charts/local/addons/${local.addon_standard.Feature}/values.yaml"
    project                                = "default"
    destination_server                     = "https://kubernetes.default.svc"
    destination_namespace                  = local.addon_standard.Feature
    syncPolicy_automated_prune             = true
    syncPolicy_automated_selfHeal          = true
    syncPolicy_syncOptions_CreateNamespace = true
  }
}

resource "kubernetes_secret_v1" "redis" {
  metadata {
    name      = "redis"
    namespace = local.addon_standard.Feature
  }

  data = {
    REDIS_PASSWORD = jsondecode(data.aws_secretsmanager_secret_version.redis.secret_string)["REDIS_PASSWORD"]
  }
  depends_on = [module.argocd_app]
}