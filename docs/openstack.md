# Provisioning and Deploying Kubernetes on OpenStack

## Step 1

Export your OpenStack CA CERT file and Initialize the kubespray repo:

```shell
export OS_CACERT=/tk8/openstack/ca.crt

tk8 cluster init
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

