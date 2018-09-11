variable "number_of_k8s_masters" {}

variable "cluster_name" {}

variable "network_id" {}

variable "floatingip_pool" {}

variable "k8s_master_priv_ips" {
  type = "list"
}

variable "k8s_master_security_group" {}
