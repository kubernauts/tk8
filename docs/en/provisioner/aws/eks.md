# Provisioning with EKS


Provide the AWS credentials in following ways:

* [Environment Variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html). You will need to specify `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.

## Create a cluster

Adapt the `config.yaml` file to specify the cluster details. [Example config](../../../../config.yaml.example):

```yaml
eks:
  cluster-name: "kubernauts-eks"
  aws_region: "us-west-2"
  node-instance-type: "m4.large"
  desired-capacity: 1
  autoscalling-max-size: 2
  autoscalling-min-size: 1  
  key-file-path: "~/.ssh/id_rsa.pub"
```

### Prerequisites using EKS

* [Git](https://git-scm.com/)
* [Terraform](https://www.terraform.io/downloads.html)
* [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Python](https://www.python.org/downloads/)
* [pip](https://pip.pypa.io/en/stable/installing/)
* [AWS IAM Authenticator](https://github.com/kubernetes-sigs/aws-iam-authenticator)
* Existing SSH keypair in AWS
* AWS access and secret keys
* Exported AWS Credentials

Once done run:

```shell
tk8ctl cluster install eks
```

or with Docker

```shell
docker run -v <path-to-the-AWS-SSH-key>:/root/.ssh/ -v "$(pwd)":/tk8 -e AWS_ACCESS_KEY_ID=xxx -e AWS_SECRET_ACCESS_KEY=XXX kubernautslabs/tk8 cluster install eks
```

Post installation the **kubeconfig** will be available at: _$(pwd)/inventory/*yourWorkspaceOrClusterName*/provisioner/kubeconfig_

**Do not** delete the inventory directory post installation as the cluster state will be saved in it.

---

## Destroy the provisioned cluster

Make sure you are in same directory where you executed `tk8 cluster install eks` with the inventory directory.
If you use a different workspace name with the --name flag please provide it on destroying too.

To delete the provisioned cluster run:

```shell
tk8ctl cluster destroy eks
```

or with Docker

```shell
docker run -v <path-to-the-AWS-SSH-key>:/root/.ssh/ -v "$(pwd)":/tk8 -e AWS_ACCESS_KEY_ID=xxx -e AWS_SECRET_ACCESS_KEY=XXX kubernautslabs/tk8 cluster destroy eks
```
