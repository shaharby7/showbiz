# Root terragrunt config
# Common configuration inherited by all environments

locals {
  env = basename(get_terragrunt_dir())
}

# Generate provider config for Kubernetes/Helm
generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}
EOF
}

# Use local backend for state (override in staging/prod for remote backend)
remote_state {
  backend = "local"
  config = {
    path = "${get_terragrunt_dir()}/terraform.tfstate"
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}
