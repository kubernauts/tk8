# Using RKE (Rancher Kubernetes Engine) to provision kubernetes cluster with TK8

Provide the AWS credentials in following ways:

* [Environment Variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html). You will need to specify `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`.

## Prerequisites

* [Git](https://git-scm.com/)
* [Terraform](https://www.terraform.io/downloads.html)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [rke-cli](https://github.com/rancher/rke)
* [terraform-provider-rke](https://github.com/yamamoto-febc/terraform-provider-rke)
* Familiarity with rke configuration.

## Usage

### Install

* Modify values in `config.yaml.example` inside `rke` key. Supported options as of now:

```plain
cluster_name - name of cluster
node_os - operating system for nodes. Currently, ubuntu
rke_aws_region - AWS region to install the cluster in.
authorization - authorization in cluster. Possible values: `rbac`, `none`. Recommended is `rbac`
rke_node_instance_type - instance type for nodes. t2.micro is highly discouraged.
node_count - number of nodes to launch in cluster
cloud_provisioner - cloud provisioner for rke cluster. Currently, only aws is supported.
```

Example:

```yaml
rke:
  cluster_name: "rke-tk8"
  node_os: ubuntu
  rke_aws_region: us-east-2
  authorization: "rbac"
  rke_node_instance_type: "t2.medium"
  node_count: 5
  cloud_provider: aws
```

* After appropriate values are filled, run:

```shell
tk8ctl cluster install rke
```

Once the infrastructure and cluster is setup, `kubeconfig` file and `rancher-cluster.yml` file will be available at:

1. kubeconfig - `inventory/rke-tk8/provisioner/kube_config_cluster.yml`
2. rancher config - `inventory/rke-tk8/provisioner/rancher-cluster.yml`

* Now you can use the same `kubeconfig` file to interact with kubernetes with kubectl and use `rancher-cluster.yml` to interact with cluster using `rke` cli.

> **Note**: Do not `mv` these files from this directory to somewhere else as these are stored in Terraform states. If required, make a copy.

### Remove

* For removing the rke cluster and keep the underlying infrastructure, run:

```shell
tk8ctl cluster remove rke
```

This is equivalent to `rke remove --config rancher-cluster.yml`.

### Destroy

* This will destroy the complete infrastructure. Run:

```shell
tk8ctl cluster destroy rke
```

> **Note** This is just a cluster provisioner, it will not install `rancher-2.x` on the cluster by itself. Use

```shell
tk8ctl addon install rancher
```

on the cluster.
