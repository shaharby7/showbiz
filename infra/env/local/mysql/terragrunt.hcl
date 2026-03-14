include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "${get_repo_root()}/infra/modules/local/mysql"
}

inputs = {
  namespace     = "showbiz"
  root_password = "rootpassword"
  database      = "showbiz"
  user          = "showbiz"
  password      = "showbiz_dev"
  storage_size  = "2Gi"
}
