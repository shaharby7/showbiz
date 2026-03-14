include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "${get_repo_root()}/infra/modules/local/minikube"
}

inputs = {
  cluster_name       = "showbiz"
  driver             = "docker"
  cpus               = 4
  memory             = 8192
  kubernetes_version = "v1.28.0"
  addons             = ["ingress", "metrics-server"]
}
