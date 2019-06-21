# Overview

The aim of this guide is to provide a reference for deploying kubernetes on a bare metal infrastructure with the option of a Load balancer service type.

This guide will make use of the kubernetes agnostic tool TK8.

TK8 is a tool that is able to install vanilla kubernetes on cloud platforms and bare metal infrastructures, it is currently based on the kubespray project.

## Prerequisites

Since this is a bare metal infrastructure, this requires that the servers must be deployed ahead of time.

Necessary informations like IP addresses and Hostnames should be retrieved ahead of time.

Installing/configuring the operating system on the servers is beyond the scope of this reference guide

Basic linux administration skill is assumed.

The following softwares are needed to have a successful deployment, these must be installed on the Deployment host:

* Ansible
* Python-netaddr
* Jinja2

## Sample Reference Architecture

![baremetal-arch](/docs/images/baremetal-arch.png)

Description for some of the appliances is below:

* Master Node: This is where the API, scheduler, controllers services is installed. An external or internal LB can be used for HA purposes, an external is recommended for optimum performance because the internal LB is deployed as a separate POD inside the kubernetes cluster, this means the LB is dependent on the kubernetes deployment and it is load balancing at the same time hence the reason to consider separation of concern. The external LB can be a seperate haproxy, nginx or hardware based like F5.

* ETCD: This is the key/value store service for the kubernetes deployment, it is strongly recommended for this to be HA with a odd number so as to achieve quorum.

* Worker Node:  This is where actual application containers are deployed, kubernetes services like the kubelet and kubeproxy resides on this node. The number of nodes is actually dependent on how much application workload will be deployed, additional nodes can be added to increase capacity.

* Storage: This is optional, that is why a dotted line is used because sometimes worker nodes utilize local storages attached to them to create the needed persistent disks for the application containers. The recommended solution is to use centralized storage solution like NFS, Ceph, Glusterfs etc. This kind of abstracted storage can be independently optimized and expanded without affecting the running application containers, also if PODs with persistent disks are to redeployed on another node, the information persists because it is not local to the worker node. Host paths should only be used to testing or development purposes.

* Deployment Node: This is where the installation will be initiated from. Ansible and it’s required modules (python-netaddr, Jinja2) will have to be installed on this node.

* Public IP Address Range: This is the range of public IPs that will be configured for use with MetalLB, it is important to note that each worker node should have a network card that is configured with an IP address from this range also because this is where the gateway settings will be reside and the LB service type will required the worker node’s routing table to properly send traffic to Edge router or Firewall.

## Getting the infrastructre ready for provisioning

* Install and prepare the servers for master, etcd, worker
* Configure key based SSH access from the Ansible host to the servers
* Create a `hosts.ini` file to be used. This is will consist details like IP address to hostname mappings, specifying the groups with the respective nodes that will be used, bastion host (if specified) etc. A sample `hosts.ini` is given below:

```ini
# ## Configure 'ip' variable to bind kubernetes services on a
# ## different ip than the default iface
node1 ansible_host=95.54.0.12  # ip=10.3.0.1
node2 ansible_host=95.54.0.13  # ip=10.3.0.2
node3 ansible_host=95.54.0.14  # ip=10.3.0.3
node4 ansible_host=95.54.0.15  # ip=10.3.0.4
node5 ansible_host=95.54.0.16  # ip=10.3.0.5
node6 ansible_host=95.54.0.17  # ip=10.3.0.6

# ## configure a bastion host if your nodes are not directly reachable
# bastion ansible_host=x.x.x.x ansible_user=some_user

[kube-master]
node1
node2

[etcd]
node1
node2
node3

[kube-node]
node2
node3
node4
node5
node6


[k8s-cluster:children]
kube-master
kube-node
```

The `ansible_host` parameter is the IP address that will be used to access the servers via SSH while the `ip` parameter is the address (recommended to be a private IP) that the kubernetes service will be listening on. Setting these two parameters is highly recommended especially if your server has multiple IP addresses or you don’t have a working DNS.

Create the inventory file on the Ansible host.

## Deploying Kubernetes

Initialize the kubespray repo:

```shell
tk8 cluster init
```

Set the neecessary parameters

Copy and replace (or edit) the sample `host.ini` file in `tk8/baremetal/`

Set the SSH username that will be used to access the servers in `tk8/baremetal/variables.yml` and accordingly set the 'become' parameter to 'no' if you are using the 'root' user or else change it to 'yes'. Sample is given below:

```yaml
os:
  username: centos
  become: "yes"
```

**N.B*:* Using root user is not recommended, instead use a separate username with sudo rights.

## Kubernetes Installation

Start the kubernetes installation on the bare metal servers:

```shell
tk8 cluster baremetal --install
```

## Configure and Deploy MetalLB

The MetalLB configmap file is also located at `tk8/baremetal/lb-config.yml`, kindly edit this file with the `public/routable IPs` that will be used for the LB service type in the kubernetes deployment. Example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
 namespace: metallb-system
 name: config
data:
 config: |
   address-pools:
   - name: my-ip-space
     protocol: layer2
     addresses:
     - 192.168.200.20-192.168.200.30
```

The address field above is the IP address range that will be used for the LB service type, the MetalLB controller POD will listen to the LB service type calls from API.

To deploy MetalLB run:

```shell
tk8 cluster baremetal --loadbalancer
```

## Post Installation

Below is the snapshot of the namespaces after installation:

```shell
kubectl get namespace
NAME             STATUS    AGE
default          Active    16d
kube-public      Active    16d
kube-system      Active    16d
metallb-system   Active    15d
```

Metallb-system namespace is created during installation, the PODs and configmap in this namespace is given below:

```shell
kubectl -n metallb-system get pods
NAME                          READY     STATUS    RESTARTS   AGE
controller-786bcb6c46-769sv   1/1       Running   1          15d
speaker-nllw7                 1/1       Running   1          15d
```

```shell
kubectl -n metallb-system describe configmap config
Name:         config
Namespace:    metallb-system
Labels:       <none>
Annotations:  <none>

Data
====
config:
----
address-pools:
- name: my-ip-space
  protocol: layer2
  addresses:
  - 192.168.200.20-192.168.200.30
```
