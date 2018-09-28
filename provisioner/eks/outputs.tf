output "kubeconfig" {
  value = "${module.eks.kubeconfig}"
}

output "config-map" {
  value = "${module.eks.config-map-aws-auth}"
}
