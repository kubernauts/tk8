resource "openstack_lb_loadbalancer_v2" "k8s" {
  name           = "${var.cluster_name}-lb"
  admin_state_up = "true"
  vip_subnet_id  = "${var.network_id}"

  security_group_ids = [
    "${var.k8s_master_security_group}",
  ]
}

resource "openstack_networking_floatingip_v2" "k8s_vip" {
  pool    = "${var.floatingip_pool}"
  port_id = "${openstack_lb_loadbalancer_v2.k8s.vip_port_id}"

  depends_on = [
    "openstack_lb_loadbalancer_v2.k8s",
  ]
}

resource "openstack_lb_listener_v2" "k8s" {
  protocol        = "HTTPS"
  protocol_port   = 6443
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.k8s.id}"
}

resource "openstack_lb_pool_v2" "k8s" {
  protocol    = "HTTPS"
  lb_method   = "ROUND_ROBIN"
  listener_id = "${openstack_lb_listener_v2.k8s.id}"
}

resource "openstack_lb_member_v2" "k8s" {
  name          = "${var.cluster_name}-k8s-lb-member"
  count         = "${var.number_of_k8s_masters}"
  pool_id       = "${openstack_lb_pool_v2.k8s.id}"
  subnet_id     = "${var.network_id}"
  address       = "${element(var.k8s_master_priv_ips, count.index)}"
  protocol_port = 6443

  depends_on = [
    "openstack_lb_loadbalancer_v2.k8s",
  ]
}

resource "openstack_lb_monitor_v2" "k8s" {
  pool_id        = "${openstack_lb_pool_v2.k8s.id}"
  type           = "HTTPS"
  delay          = 20
  timeout        = 10
  max_retries    = 5
  url_path       = "/"
  expected_codes = 200
}
