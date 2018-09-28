output "loadbalancer_id" {
  value = "${openstack_lb_loadbalancer_v2.k8s.vip_port_id}"
}

output "k8s_master_vip" {
  value = "${openstack_networking_floatingip_v2.k8s_vip.address}"
}
