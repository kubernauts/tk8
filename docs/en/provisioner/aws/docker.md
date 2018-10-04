# Provisioning with Docker

## Using Docker image

Create a new Folder with a config file and the ssh keys and mount it to the docker container

## Create a cluster

Adapt the `config.yaml` file to specify the cluster details. [Example config](https://raw.githubusercontent.com/kubernauts/tk8/master/config.yaml.example):

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
vi config.yaml --> provide your AWS access and secret key and an exsiting SSH keypair in AWS
docker run -v ~/.ssh/:/root/.ssh/ -v "$(pwd)":/tk8 kubernautslabs/tk8 cluster install aws
```

Post installation the **kubeconfig** will be available at: _$(pwd)/inventory/*yourWorkspaceOrClusterName*/artifacts/admin.conf_

**Do not** delete the inventory directory post installation as the cluster state will be saved in it.

---

## Destroy the provisioned cluster

Make sure you are in same directory where you executed `tk8 cluster install aws` with the inventory directory.
If you use a different workspace name with the --name flag please provided it on destroying too.

To delete the provisioned cluster run:

```shell
docker run -v ~/.ssh/:/root/.ssh/ -v "$(pwd)":/tk8 kubernautslabs/tk8 cluster destroy aws
```