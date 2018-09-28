provider "openstack" {
  cloud = "mycloud"
}

module "network" {
  source = "modules/network"

  external_net    = "${var.external_net}"
  network_name    = "${var.network_name}"
  cluster_name    = "${var.cluster_name}"
  dns_nameservers = "${var.dns_nameservers}"
}

module "ips" {
  source = "modules/ips"

  number_of_k8s_masters         = "${var.number_of_k8s_masters}"
  number_of_k8s_masters_no_etcd = "${var.number_of_k8s_masters_no_etcd}"
  number_of_k8s_nodes           = "${var.number_of_k8s_nodes}"
  floatingip_pool               = "${var.floatingip_pool}"
  number_of_bastions            = "${var.number_of_bastions}"
  external_net                  = "${var.external_net}"
  network_name                  = "${var.network_name}"
  router_id                     = "${module.network.router_id}"
}

module "compute" {
  source = "modules/compute"

  cluster_name                                 = "${var.cluster_name}"
  number_of_k8s_masters                        = "${var.number_of_k8s_masters}"
  number_of_k8s_masters_no_etcd                = "${var.number_of_k8s_masters_no_etcd}"
  number_of_etcd                               = "${var.number_of_etcd}"
  number_of_k8s_masters_no_floating_ip         = "${var.number_of_k8s_masters_no_floating_ip}"
  number_of_k8s_masters_no_floating_ip_no_etcd = "${var.number_of_k8s_masters_no_floating_ip_no_etcd}"
  number_of_k8s_nodes                          = "${var.number_of_k8s_nodes}"
  number_of_bastions                           = "${var.number_of_bastions}"
  number_of_k8s_nodes_no_floating_ip           = "${var.number_of_k8s_nodes_no_floating_ip}"
  number_of_gfs_nodes_no_floating_ip           = "${var.number_of_gfs_nodes_no_floating_ip}"
  gfs_volume_size_in_gb                        = "${var.gfs_volume_size_in_gb}"
  public_key_path                              = "${var.public_key_path}"
  image                                        = "${var.image}"
  image_gfs                                    = "${var.image_gfs}"
  ssh_user                                     = "${var.ssh_user}"
  ssh_user_gfs                                 = "${var.ssh_user_gfs}"
  flavor_k8s_master                            = "${var.flavor_k8s_master}"
  flavor_k8s_node                              = "${var.flavor_k8s_node}"
  flavor_etcd                                  = "${var.flavor_etcd}"
  flavor_gfs_node                              = "${var.flavor_gfs_node}"
  network_name                                 = "${var.network_name}"
  flavor_bastion                               = "${var.flavor_bastion}"
  k8s_node_fips                                = "${module.ips.k8s_node_fips}"
  bastion_fips                                 = "${module.ips.bastion_fips}"

  network_id      = "${module.network.router_id}"
  loadbalancer_id = "${module.lb.loadbalancer_id}"
}

module "lb" {
  source = "modules/lb"

  cluster_name              = "${var.cluster_name}"
  number_of_k8s_masters     = "${var.number_of_k8s_masters}"
  network_id                = "${module.network.network_id}"
  k8s_master_priv_ips       = "${module.compute.k8s_master_priv_ips}"
  floatingip_pool           = "${var.floatingip_pool}"
  k8s_master_security_group = "${module.compute.k8s_master_security_group}"
}

output "private_subnet_id" {
  value = "${module.network.network_id}"
}

output "floating_network_id" {
  value = "${var.external_net}"
}

output "router_id" {
  value = "${module.network.router_id}"
}

output "k8s_master_lb_vip" {
  value = "${module.lb.k8s_master_vip}"
}

output "k8s_node_fips" {
  value = "${module.ips.k8s_node_fips}"
}

output "bastion_fips" {
  value = "${module.ips.bastion_fips}"
}

/*
* Create Kubespray Inventory File
*
*/
data "template_file" "inventory" {
  template = "${file("${path.module}/inventory.tpl")}"

  vars {
    public_ip_address_bastion = "${join("\n",formatlist("bastion ansible_host=%s" , "${module.ips.bastion_fips}"))}"
    connection_strings_master = "${join("\n",formatlist("%s ansible_host=%s","${module.compute.k8s_master_names}", "${module.compute.k8s_master_priv_ips}"))}"
    connection_strings_node   = "${join("\n",formatlist("%s ansible_host=%s","${module.compute.k8s_worker_names}", "${module.compute.k8s_worker_priv_ips}"))}"
    connection_strings_etcd   = "${join("\n",formatlist("%s ansible_host=%s","${module.compute.k8s_etcd_names}", "${module.compute.k8s_etcd_priv_ips}"))}"
    list_master               = "${join("\n","${module.compute.k8s_master_names}")}"
    list_node                 = "${join("\n","${module.compute.k8s_worker_names}")}"
    list_etcd                 = "${join("\n","${module.compute.k8s_etcd_names}")}"
    elb_api_fqdn              = "apiserver_loadbalancer_domain_name=\"${var.elb_api_fqdn}\""
  }
}

resource "null_resource" "inventories" {
  provisioner "local-exec" {
    command = "echo '${data.template_file.inventory.rendered}' > ./hosts.ini"
  }

  triggers {
    template = "${data.template_file.inventory.rendered}"
  }
}

resource "null_resource" "export-floating-network-id" {
  provisioner "local-exec" {
    command = "echo 'lbaas-floating-network-id: ''${var.external_net}' > ./network-config.yaml"
  }
}

resource "null_resource" "export-priv-subnet-id" {
  provisioner "local-exec" {
    command = "echo 'lbaas-private-subnet-id: ''${module.network.network_id}' >> ./network-config.yaml"
  }
}

resource "null_resource" "export-floating-vip-master" {
  provisioner "local-exec" {
    command = "echo 'floating-master-lb-vip: ''${module.lb.k8s_master_vip}' >> ./network-config.yaml"
  }
}
