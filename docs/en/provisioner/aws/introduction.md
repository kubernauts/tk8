# Provisioning and Deploying Kubernetes on AWS

TK8 offers two possibilities to create a cluster on the desired platform. On one hand, there is the path via the CLI, which requires dependencies on the executing system. Alternatively, it is possible to create a cluster on the desired platform with a minimal dependency on Docker.

## Provisioning using CLI

The use of the CLI gives you more freedom to individualize your workflows.

### Prerequisites using CLI

* [Git](https://git-scm.com/)
* [Terraform](https://www.terraform.io/downloads.html)
* [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Python](https://www.python.org/downloads/)
* [pip](https://pip.pypa.io/en/stable/installing/)
* Existing SSH keypair in AWS
* AWS access and secret keys

[Documentation](cli.md)

## Provisioning using Docker

This is the easiest way to deploy a H.A. Kubernetes cluster on AWS with the [kubernautslabs/tk8](https://hub.docker.com/r/kubernautslabs/tk8) docker image.

### Prerequisites using Docker

* [Git](https://git-scm.com/)
* [Docker](https://docs.docker.com/install/)
* Existing SSH keypair in AWS
* AWS access and secret keys

[Documentation](docker.md)

## Provisioning using EKS

This is the easiest way to deploy an AWS EKS cluster using the [kubernautslabs/tk8](https://hub.docker.com/r/kubernautslabs/tk8) docker image.

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

[Documentation](eks.md)
