![Logo](docs/images/tk8.png)

# Tk8 Rest API
TK8 now also provides the ability to create clusters via REST API.
The configuration stored for the created clusters will be either stored on the local machine path provided by the user
or the S3 bucket provided by the user


## Usage

By Default the REST server starts at port 8091. However that can changed to a port you desire.

server-mode provides the following flags

Flags:
  -s, --config-store string          Storage for config files - local or s3 (default "local")
  -a, --config-store-path string     Storage location for config files - directory path for local (default ".") or bucket name for s3
  -r, --config-store-region string   Region for S3 bucket
  -p, --port uint16                  Port number for the Tk8 rest api (default 8091)

To start the server using local storage, execute

```
tk8 server-mode -p 8091 -s local -a /path/to/the/local/directory/you/want
```

To start the server using s3 bucket, execute
```
tk8 server-mode -p 8091 -s s3 -a <bucket_name> -r <region_of_the_s3_bucket>
```

## Creating Clusters

To create a cluster follow the configuration params required by each provisioner and create a POST request against the provisioner
Payload to be in JSON

As an example

To create an AWS cluster, use the below template to create a request payload and save as aws-create.json

```
{
   "clustername": "tk8-aws-cluster",
   "os": "centos",
   "aws_access_key_id": "",
   "aws_secret_access_key": "",
   "aws_ssh_keypair": "",
   "aws_default_region": "",
   "aws_vpc_cidr_block": "10.250.192.0/18",
   "aws_cidr_subnets_private": "[\"10.250.192.0/20\",\"10.250.208.0/20\"]",
   "aws_cidr_subnets_public": "[\"10.250.224.0/20\",\"10.250.240.0/20\"]",
   "aws_bastion_size": "t2.medium",
   "aws_kube_master_num": 1,
   "aws_kube_master_size": "t2.medium",
   "aws_etcd_num": 1,
   "aws_etcd_size": "t2.medium",
   "aws_kube_worker_num": 1,
   "aws_kube_worker_size": "t2.medium",
   "aws_elb_api_port": 6443,
   "k8s_secure_api_port": 6443,
   "kube_insecure_apiserver_address": "0.0.0.0",
   "kubeadm_enabled": false,
   "kube_network_plugin": "flannel"
}
```

Execute a curl request to AWS endpoint

```
curl -X POST  -d @aws-create.json localhost:8091/v1/cluster/aws
```



Currently the supported provisioners are 

 - AWS provisioner
 - EKS Provisioner
 - RKE Provisioner


 ## Rest Endpoints

### Get cluster details 

 - GET /v1/cluster/aws/{id}
 - GET /v1/cluster/rke/{id}
 - GET /v1/cluster/eks/{id}

### Delete created Cluster

  - DELETE /v1/cluster/aws/{id}
  - DELETE /v1/cluster/eks/{id}
  - DELETE /v1/cluster/rke/{id}

### Create cluster
 - POST /v1/cluster/aws
 - POST /v1/cluster/rke
 - POST /v1/cluster/eks