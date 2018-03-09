![Screenshot](tk8.png)
# tk8 (The Kubernaut): A multi-cloud, multi-cluster Kubernetes (k8s) platform installation and integration tool

tk8 is a CLI written in Golang to deploy Kubernetes using terraform and kubespray and also to install additional addons such as jmeter for loadtesting on k8s.

## TL;DR

git clone https://github.com/kubernauts/tk8.git

docker build -t tk8 ./tk8/.

vi ./tk8/config.yaml

--> pls. provide the aws access and secret keys, your ssh keypair name and your desired aws region

Note: create a key pair in aws with your publich key (id_rsa.pub)

docker run -it --rm -v "$(pwd)"/tk8:/tk8 tk8 bash

mkdir .ssh

vi .ssh/id_rsa --> paste yor private key in id_rsa in the container

chmod 400 .ssh/id_rsa

tk8 cluster init ### kubespray will be cloned

tk8 cluster aws -c ### terraform brings up 2 bastion hosts and the master and worker nodes as specified in config.yaml, wail till the hosts are initialized

tk8 cluster aws -i ### kubespray installs kubernetes

