# Provisioning with CLI

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
* [Config file](https://raw.githubusercontent.com/kubernauts/tk8/master/config.yaml.example)

## Create a cluster

Adapt the `config.yaml` file to specify the cluster details. [Example config](../../../../config.yaml.example):

```yaml
aws:
   clustername: kubernauts
   os: centos
   aws_access_key_id: ""
   aws_secret_access_key: ""
   aws_ssh_keypair: ""
   aws_default_region: ""
   aws_vpc_cidr_block : "10.250.192.0/18"
   aws_cidr_subnets_private : '["10.250.192.0/20","10.250.208.0/20"]'
   aws_cidr_subnets_public : '["10.250.224.0/20","10.250.240.0/20"]'
   aws_bastion_size : "t2.medium"
   aws_kube_master_num : 1
   aws_kube_master_size : "t2.medium"
   aws_etcd_num : 1
   aws_etcd_size : "t2.medium"
   aws_kube_worker_num : 1
   aws_kube_worker_size : "t2.medium"
   aws_elb_api_port : 6443
   k8s_secure_api_port : 6443
   kube_insecure_apiserver_address : "0.0.0.0"
   kubeadm_enabled: false
   kube_network_plugin: "flannel"
```

Once done run:

```shell
tk8 cluster install aws
```

Post installation the **kubeconfig** will be available at: _$(pwd)/inventory/*yourWorkspaceOrClusterName*/artifacts/admin.conf_

**Do not** delete the inventory directory post installation as the cluster state will be saved in it.

---

## Destroy the provisioned cluster

Make sure you are in same directory where you executed `tk8 cluster install aws` with the inventory directory.
If you use a different workspace name with the --name flag please provided it on destroying too.

To delete the provisioned cluster run:

```shell
tk8 cluster destroy aws
```
