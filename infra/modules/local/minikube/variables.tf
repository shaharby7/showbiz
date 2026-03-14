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

variable "addons" {
  description = "List of Minikube addons to enable"
  type        = list(string)
  default     = ["ingress", "metrics-server"]
}
