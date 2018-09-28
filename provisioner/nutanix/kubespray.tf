provider "nutanix" {
  username = "root"
  password = "nx2Tech673!"
  endpoint = "10.21.52.37"
  insecure = true
  port     = 9440
}
module "compute" {
  source = "modules/compute"
}

# data "template_file" "inventory" {
#   template = "${file("${path.module}/inventory.tpl")}"

#   vars {
#     public_ip_address_bastion = "${join("\n",formatlist("bastion ansible_host=%s" , "${module.ips.bastion_fips}"))}"
#     connection_strings_master = "${join("\n",formatlist("%s ansible_host=%s","${module.compute.k8s_master_names}", "${module.compute.k8s_master_priv_ips}"))}"
#     connection_strings_node   = "${join("\n",formatlist("%s ansible_host=%s","${module.compute.k8s_worker_names}", "${module.compute.k8s_worker_priv_ips}"))}"
#     connection_strings_etcd   = "${join("\n",formatlist("%s ansible_host=%s","${module.compute.k8s_etcd_names}", "${module.compute.k8s_etcd_priv_ips}"))}"
#     list_master               = "${join("\n","${module.compute.k8s_master_names}")}"
#     list_node                 = "${join("\n","${module.compute.k8s_worker_names}")}"
#     list_etcd                 = "${join("\n","${module.compute.k8s_etcd_names}")}"
#     elb_api_fqdn              = "apiserver_loadbalancer_domain_name=\"${var.elb_api_fqdn}\""
#   }
# }


# resource "null_resource" "inventories" {
#   provisioner "local-exec" {
#     command = "echo '${data.template_file.inventory.rendered}' > ./hosts.ini"
#   }

#   triggers {
#     template = "${data.template_file.inventory.rendered}"
#   }
# }

# resource "null_resource" "export-floating-network-id" {
#   provisioner "local-exec" {
#     command = "echo 'lbaas-floating-network-id: ''${var.external_net}' > ./network-config.yaml"
#   }
# }

# resource "null_resource" "export-priv-subnet-id" {
#   provisioner "local-exec" {
#     command = "echo 'lbaas-private-subnet-id: ''${module.network.network_id}' >> ./network-config.yaml"
#   }
# }

# resource "null_resource" "export-floating-vip-master" {
#   provisioner "local-exec" {
#     command = "echo 'floating-master-lb-vip: ''${module.lb.k8s_master_vip}' >> ./network-config.yaml"
#   }
# }
