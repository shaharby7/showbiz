terraform {
  required_version = ">= 1.0"
}

variable "cluster_name" {
  description = "Name of the Minikube cluster"
  type        = string
  default     = "showbiz"
}

variable "driver" {
  description = "Minikube driver (docker, hyperkit, etc.)"
  type        = string
  default     = "docker"
}

variable "cpus" {
  description = "Number of CPUs for the Minikube VM"
  type        = number
  default     = 4
}

variable "memory" {
  description = "Memory in MB for the Minikube VM"
  type        = number
  default     = 8192
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "v1.28.0"
}

resource "null_resource" "minikube" {
  provisioner "local-exec" {
    command = <<-EOT
      minikube start \
        --profile=${var.cluster_name} \
        --driver=${var.driver} \
        --cpus=${var.cpus} \
        --memory=${var.memory} \
        --kubernetes-version=${var.kubernetes_version} \
        --addons=ingress,metrics-server
    EOT
  }

  provisioner "local-exec" {
    when    = destroy
    command = "minikube delete --profile=showbiz"
  }
}

output "cluster_name" {
  value = var.cluster_name
}

output "kubeconfig_path" {
  value = "~/.kube/config"
}
