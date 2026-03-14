include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "${get_repo_root()}/infra/modules/local/kubevirt"
}

dependency "minikube" {
  config_path = "../minikube"
}

inputs = {
  kubevirt_version = "v1.2.0"
  vmis_namespace   = "vmis"
}
