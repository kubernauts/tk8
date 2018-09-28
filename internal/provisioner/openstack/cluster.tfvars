# your Kubernetes cluster name here
cluster_name = "kubernauts"

# SSH key to use for access to nodes
public_key_path = "/home/chris/Documents/projects/kubernauts/openstack_admin_ssh_public"

# image to use for bastion, masters, standalone etcd instances, and nodes
image = "Centos 7"

# user on the node (ex. core on Container Linux, ubuntu on Ubuntu, etc.)
ssh_user = "centos"

# 0|1 bastion nodes
number_of_bastions = 1

flavor_bastion = "8"

# standalone etcds
number_of_etcd = 1

flavor_etcd = "7"

# masters
number_of_k8s_masters = 2

number_of_k8s_masters_no_etcd = 0

number_of_k8s_masters_no_floating_ip = 0

number_of_k8s_masters_no_floating_ip_no_etcd = 0

flavor_k8s_master = "7"

# nodes
number_of_k8s_nodes = 0

number_of_k8s_nodes_no_floating_ip = 2

flavor_k8s_node = "7"

# GlusterFS
# either 0 or more than one
#number_of_gfs_nodes_no_floating_ip = 0
#gfs_volume_size_in_gb = 150
# Container Linux does not support GlusterFS
#image_gfs = "<image name>"
# May be different from other nodes
#ssh_user_gfs = "ubuntu"
#flavor_gfs_node = "<UUID>"

# networking
network_name = "kubernauts"

external_net = "b1152444-f5ef-4c3d-a4b5-73324bfae4be"

floatingip_pool = "external"

elb_api_fqdn = "tk8.example.com"
