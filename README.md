# TK8: A multi-cloud, multi-cluster Kubernetes platform installation and integration tool

TK8 is a command line tool written in Go and fully automatically installs Kubernet on any environment. With TK8 you are able to centrally manage different Kubernet clusters with different configurations. In addition, TK8 with its simple add-on integration offers the possibility to quickly, cleanly and easily distribute extensions to the different Kubernetes clusters.

These include a Jmeter cluster for load testing, Prometheus for monitoring, Jaeger, Linkerd or Zippkin for tracing, Ambassador API Gateway with Envoy for Ingress and Load Balancing, Istio as mesh support solution, Jenkins-X for CI/CD integration. In addition, the add-on system also supports the management of helm packages.

## Table of content

The documentation as well as a detailed table of contents can be found here.

* [Table of content](docs/en/SUMMARY.md)

## Installation

The TK8 CLI requires some dependencies to perform its tasks.
At the moment we still need your help here, but we are already working on a setup script that will do these tasks for you.

### Terraform

Terraform is required to automatically set up the infrastructure in the desired environment.
[Terraform Installation](https://www.terraform.io/intro/getting-started/install.html)

### Ansible

Ansible is required to run the automated installation routines in the desired and automatically created environment.
[Ansible Installation](https://docs.ansible.com/ansible/2.5/installation_guide/intro_installation.html#installing-the-control-machine)

### Kubectl

Kubectl is needed by the CLI to roll out the add-ons and by you to access your clusters.
[kubectl Installation](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Python and pip

In the automated routines Python scripts are used, in addition with Pip dependencies are loaded for this.
[Python Installation](https://www.python.org/downloads/)
[pip Installation](https://pip.pypa.io/en/stable/installing/)

### AWS IAM Authenticator

If you want to install an EKS cluster with TK8, the [AWS IAM Authenticator](https://github.com/kubernetes-sigs/aws-iam-authenticator) must be executable _(/usr/local/bin)_. This is included in the provisioner package EKS of the TK8 CLI or can be found in the given link.

## Usage

Since there are different target platforms with the TK8 CLI and we have described these separately in detail in the documentation, we would like to give you just one example using AWS.

Download the executable file for your operating system from the Release section or build your own version with the `go build` command.

Create a separate folder and store the executable file there, a configuration file is also required. This file can be found under the name config.yaml.example. Enter here the necessary parameters for your cluster as well as the AWS CLI Key and the Secret. Additionally you should put your AWS credentials in the environment variables because parts of the CLI (EKS cluster) need them there.

`export AWS_SECRET_ACCESS_KEY=xxx`
`export AWS_ACCESS_KEY_ID=xxx`

They then execute the CLI with the command:
`tk8 cluster install aws`

With this command the TK8 CLI creates all required resources in AWS and installs a Kubernet cluster for it.

If you no longer need the cluster, you can use the command:
`tk8 cluster destroy aws`
to automatically remove all resources.

## Contributing

For the provision of add-ons we have a separate documentation area and examples how you can build your extensions and integrate them into the TK8 project. You can also reach us at Slack.

As a platform provider we have a separate documentation area here which is only about integrating a platform in TK8. Here you will find detailed instructions and examples on how TK8 will execute your integration. You can also reach us in slack.

To participate in the core, please create an issue or get in touch with us in Slack.

Get in touch
[Join us on Kubernauts Slack Channel](https://kubernauts-slack-join.herokuapp.com/)

## Credits

Founder and initiator of this project is [Arash Kaffamanesh](https://github.com/arashkaffamanesh) Founder and CEO of [cloudssky GmbH](https://cloudssky.com/de/) and [Kubernauts GmbH](https://kubernauts.de/en/home/)

The project is supported by cloud computing experts from cloudssky GmbH and Kubernauts GmbH.
[Christopher Adigun](https://github.com/infinitydon)
[Arush Salil](https://github.com/arush-sal)
[Manuel Müller](https://github.com/MuellerMH)
[Niki](https://github.com/niki-1905)
[Anoop](https://github.com/anoopl)

A big thanks goes to the contributors of [Kubespray](https://github.com/kubernetes-incubator/kubespray) whose great work we use as a basis for the setup and installation of Kubernetes in the AWS Cloud.

Furthermore we would like to thank the contributors of [kubeadm](https://github.com/kubernetes/kubernetes/tree/master/cmd/kubeadm) which is currently not only part of the Kubespray project, but also of the TK8.

Ebenfalls ein großes Dankeschön an [Wesley Charles Blake] (https://github.com/WesleyCharlesBlake), auf dessen Grundlage wir unsere EKS-Integration anbieten konnten.

## Lizenz

[Tk8 Apache Lizenz](LIZENZ)
[MIT Lizenz EKS](https://github.com/kubernauts/tk8eks/blob/master/LICENSE-Wesley-Charles-Blake)
[MIT Lizenz EKS](https://github.com/kubernauts/tk8eks/blob/master/LICENSE)
