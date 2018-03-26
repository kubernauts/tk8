output "k8s_master_priv_ips" {
  value = ["${openstack_compute_instance_v2.k8s_master.*.network.0.fixed_ip_v4}"]
}

output "k8s_worker_priv_ips" {
  value = ["${openstack_compute_instance_v2.k8s_node_no_floating_ip.*.network.0.fixed_ip_v4}"]
}

output "k8s_etcd_priv_ips" {
  value = ["${openstack_compute_instance_v2.etcd.*.network.0.fixed_ip_v4}"]
}

output "k8s_master_names" {
  value = ["${openstack_compute_instance_v2.k8s_master.*.name}"]
}

output "k8s_worker_names" {
  value = ["${openstack_compute_instance_v2.k8s_node_no_floating_ip.*.name}"]
}

output "k8s_etcd_names" {
  value = ["${openstack_compute_instance_v2.etcd.*.name}"]
}

output "k8s_master_security_group" {
  value = "${openstack_compute_secgroup_v2.k8s_master.id}"
}
