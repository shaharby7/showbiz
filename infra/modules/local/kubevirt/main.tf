resource "kubernetes_namespace" "vmis" {
  metadata {
    name = var.vmis_namespace
  }
}

resource "kubernetes_namespace" "kubevirt" {
  metadata {
    name = "kubevirt"
  }
}

resource "null_resource" "kubevirt_operator" {
  depends_on = [kubernetes_namespace.kubevirt]

  provisioner "local-exec" {
    command = "kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${var.kubevirt_version}/kubevirt-operator.yaml"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "kubectl delete -f https://github.com/kubevirt/kubevirt/releases/download/v1.2.0/kubevirt-operator.yaml --ignore-not-found"
  }
}

resource "null_resource" "kubevirt_cr" {
  depends_on = [null_resource.kubevirt_operator]

  provisioner "local-exec" {
    command = "kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${var.kubevirt_version}/kubevirt-cr.yaml"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "kubectl delete -f https://github.com/kubevirt/kubevirt/releases/download/v1.2.0/kubevirt-cr.yaml --ignore-not-found"
  }
}

resource "null_resource" "wait_for_kubevirt" {
  depends_on = [null_resource.kubevirt_cr]

  provisioner "local-exec" {
    command = "kubectl -n kubevirt wait kubevirt kubevirt --for condition=Available --timeout=900s"
  }
}
