![Logo](docs/images/tk8.png)

# TK8: A multi-cloud, multi-cluster Kubernetes platform installation and integration tool based on Kubespray

TK8 is a CLI written in Golang to deploy the upstream Vanilla Kubernetes fully automated based on [Kubespray](https://github.com/kubernetes-incubator/kubespray) on any environment. **We'll provide kubeadm support as soon kubeadm HA support is available through Kubespray project.**

With TK8 you can also install additional add-ons such as a Jmeter Cluster for load testing, Prometheus for monitoring, Jaeger, Linkerd or Zippkin for tracing, Ambassador API Gateway with Envoy for ingress and load balancing, Istio for service mesh support , Jenkins-X for CI/CD and Helm or Kedge for packaging on Kubernetes.

This release supports the Kubernetes installation on AWS and OpenStack / Bare-Metal with HA capabilities.

N.B: MS Azure and GCP support will be added in the very near future.

## Installation

### Building from source

```shell
go get -u github.com/kubernauts/tk8
```


## Usage

You can either use the cli to install a supported addon in an existing Kubernetes cluster or to provision and install Kubernetes on the supported platforms. The basic usage instructions are as below:

```shell
Usage:
  tk8 [command]

Available Commands:
  addon       Install kubernetes addon packages
  cluster     Used to create kubernetes cluster on cloud infrastructures
  help        Help about any command

Flags:
      --config string   uses the config.yaml
  -h, --help            help for tk8
  -t, --toggle          Help message for toggle

Use "tk8 [command] --help" for more information about a command.
```

### Install a supported addon in an existing Kubernetes cluster

#### Using installed binary

```shell
Usage:
  tk8 addon [flags]

Flags:
  -m, --heapster   Deploy Heapster
  -h, --help       help for addon
  -l, --ltaas      Deploy Load Testing As A Service
  -p, --prom       Deploy prometheus

Global Flags:
      --config string   Path to the config.yaml
```


### Provision and install Kubernetes on the supported platforms

#### Using installed binary

```shell
Usage:
  tk8 cluster [flags]
  tk8 cluster [command]

Available Commands:
  aws         Manages the infrastructure on AWS
  baremetal   Manages the infrastructure on Baremetal
  init        Initialise kubespray repository
  openstack   Manages the infrastructure on Openstack

Flags:
  -h, --help   help for cluster

Global Flags:
      --config string   Path to the config.yaml

Use "tk8 cluster [command] --help" for more information about a command.
```

Specific platform instructions can be found in the [official documentation](https://kubernauts.gitbooks.io/tk8/content/) or in [docs](docs/)

