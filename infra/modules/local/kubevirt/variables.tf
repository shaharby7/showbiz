variable "kubevirt_version" {
  description = "KubeVirt release version"
  type        = string
  default     = "v1.2.0"
}

variable "vmis_namespace" {
  description = "Namespace for virtual machine instances"
  type        = string
  default     = "vmis"
}
