# Provisioning and Deploying Kubernetes on AWS

## Using Docker image

This is the easiest way to deploy an HA Kubernetes cluster on AWS with kubernautslabs/tk8 docker image.

## Prerequisites

* Git
* Docker
* Existing SSH keypair in AWS
* AWS access and secret keys

```shell
git clone https://github.com/kubernauts/tk8
cd tk8
docker build -t kubernautslabs/tk8 . # optional, only if you'd like to build your own docker image
vi config.yaml --> provide your AWS access and secret key and an exsiting SSH keypair in AWS
docker run -it --name tk8 -v ~/.ssh/:/root/.ssh/ -v "$(pwd)":/tk8 kubernautslabs/tk8 sh
cd tk8
tk8 cluster init # clone kubernauts/kubespray
tk8 cluster aws --create # create the cluster infra using terraform
pip install -r kubespray/requirements.txt # not needed with kubernauts/kubespray
tk8 cluster aws --install
KUBECONFIG=./kubespray/inventory/awscluster/artifacts/admin.conf kubectl get pods --all-namespaces
tk8 cluster aws --destroy
```

---

## Provision the cluster using CLI

## Prerequisites

* [Git](https://git-scm.com/)
* [Terraform](https://www.terraform.io/downloads.html)
* [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)

---

You can download the TK8 CLI for Linux and Mac OSX under releases [here](https://github.com/kubernauts/tk8/releases).

Make sure you have the sufficient permissions to create following resources in AWS

* VPC with Public/Private Subnets, and NAT Gateways in different Availability Zones
* EC2 instances used for bastion host, masters, etcd, and worker nodes
* AWS ELB in the Public Subnet for accessing the Kubernetes API from the internet
* IAM Roles which will be used with the nodes

Provide the AWS credentials in either of the following ways:

* [AWS IAM Instance Profile](https://docs.aws.amazon.com/codedeploy/latest/userguide/getting-started-create-iam-instance-profile.html)
* [Configuration and Credential Files](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html)
* [Environment Variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html). You will need to specify `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.
* [Config file](../config.yaml)

## Clone the tk8 repo

```shell
git clone https://github.com/kubernauts/tk8
cd tk8
```

Adapt the `config.yaml` file to specify the cluster details. [Example config](../config.yaml):

```yaml
aws:
   clustername: kubernauts
   os: centos # valid options are coreos/ubuntu/centos
   aws_access_key_id: # optional, see above.
   aws_secret_access_key: # optional, see above.
   aws_ssh_keypair: # needs to be an existing SSH keypair in AWS
   aws_default_region:
   aws_vpc_cidr_block : "10.250.192.0/18"
   aws_cidr_subnets_private : '["10.250.192.0/20","10.250.208.0/20"]'
   aws_cidr_subnets_public : '["10.250.224.0/20","10.250.240.0/20"]'
   aws_bastion_size : "t2.medium"
   aws_kube_master_num : 1
   aws_kube_master_size : "t2.medium"
   aws_etcd_num : 1
   aws_etcd_size : "t2.medium"
   aws_kube_worker_num : 2
   aws_kube_worker_size : "t2.medium"
   aws_elb_api_port : 6443
   k8s_secure_api_port : 6443
   kube_insecure_apiserver_address : "0.0.0.0"
```

Once done run:

```shell
tk8 cluster init
tk8 cluster aws --create
tk8 cluster aws --install
```

Post installation the **kubeconfig** will be available at: _./kubespray/inventory/awscluster/artifacts/admin.conf_

**Do not** delete the kubespray directory post installation as the cluster state will be saved in it.

---

## Destroy the provisioned cluster

Make sure you are in same directory where you executed `tk8 cluster init` with the cloned kubespray directory.

To delete the provisioned cluster run:

```shell
tk8 cluster aws --destroy
```

