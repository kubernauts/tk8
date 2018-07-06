# Provisioning and Deploying Kubernetes on OpenStack

## Prerequisites

* [Git](https://git-scm.com/)
* [Terraform](https://www.terraform.io/downloads.html)
* [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)

## Step 1

Get the tk8 repo

```shell
git clone https://github.com/kubernauts/tk8
cd tk8
wget https://github.com/kubernauts/tk8/releases/download/0.3/tk8-linux-opentack-amd64
chmod +x tk8-linux-opentack-amd64
mv tk8-linux-opentack-amd64 /user/local/bin/tk8
```

Source your OpenStack rc file , export your OpenStack CA CERT file and Initialize the kubespray repo:

```shell
source project-openrc.sh

export OS_CACERT=/path-to/tk8/openstack/ca.crt # needed if you use a self signed certificate

tk8 cluster init # initialize the kubespray repo
```

## Step 2

Adjust the following openstack files with your specific settings:

_**~/tk8/openstack/cluster.tfvars**_ -- This file contains the parameters for the cloud components to be deployed e.g flavor, image, cluster\_name, external network id \(float network\), number of masters, number etcd nodes, number of worker nodes etc. Some values have been set but some needs to be modified to suit your environment like the image, flavor \(this is the ID not name\) and external network id.

The elb\_api\_fqdn domain name should be set to a value that you can resolve within your network/deployment, this will should be resolvable to the load balancer VIP address of the master nodes i.e. Floating IP address \(you can check the load balancer on Horizon or openstack CLI to check the value of the Load Balancer floating IP address for the master nodes\).

N.B -- Currently you can only use Centos 7 image \(image can be downloaded from the centos public repo\), Ubuntu and Coreos will be added soon.

Also make sure you use the right set of SSH keys, you are to use the SSH public key that you used to create the key pair on Openstack, otherwise the kubernetes installation will not be able to access the VMs via ssh.

_**~/tk8/openstack/stack\_credentials.yaml**_  -- This file contains your openstack credentials, modify as per your openstack deployment.

## Step 3

Create the infrastructure with the following command:

```shell
tk8 cluster openstack --create
```

### Step 4

Install kubernetes on the infrastructure that was created with the following command:

```shell
tk8 cluster openstack --install
```

For the Openstack deployment, you can get the kubeconfig by checking the following location after the kubernetes installation finishes: _**./kubespray/inventory/stackcluster/artifacts/admin.conf**_

### Step 5

Destroy the infrastructure, use the following command:

```shell
tk8 cluster openstack --destroy
```

N.B -- Before destroying the cluster, make sure you delete any load balancer that was created inside your kubernetes cluster, otherwise, terraform destroy will not work completely since terraform did not create the additional load balancer \(the load balancer is tied to some other aspects of the cloud which will affect the terraform destroy procedure\).

## Using Docker image

**No Prerequisites**, oh yes you need Docker :-\)

```shell
git clone https://github.com/kubernauts/tk8
cd tk8
vi openstack/cluster.tfvars
vi openstack/stack-credentials.yaml
docker run -it --name tk8 -v ~/.ssh/:/root/.ssh/ -v "$(pwd)":/tk8 kubernautslabs/tk8 sh
cd tk8
tk8 cluster init
tk8 cluster openstack --create
pip install -r kubespray/requirements.txt
tk8 cluster openstack --install
exit
KUBECONFIG=./kubespray/inventory/awscluster/artifacts/admin.conf kubectl get pods --all-namespaces
```



