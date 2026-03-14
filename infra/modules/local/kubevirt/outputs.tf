output "vmis_namespace" {
  description = "Namespace where VMs are created"
  value       = var.vmis_namespace
}

output "kubevirt_version" {
  description = "Installed KubeVirt version"
  value       = var.kubevirt_version
}
