include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "${get_repo_root()}/infra/modules/k8s/argocd"
}

dependency "minikube" {
  config_path = "../minikube"
}

inputs = {
  namespace              = "argocd"
  chart_version          = "5.51.6"
  environment            = "local"
  repo_url               = "https://github.com/showbiz-io/showbiz.git"
  target_revision        = "main"
  app_of_apps_chart_path = "${get_repo_root()}/helm/charts/app-of-apps"
}
