package templates

var Terraform = `
aws_cluster_name = "{{.AwsClusterName}}"
aws_vpc_cidr_block = "{{.AwsVpcCidrBlock}}"
aws_cidr_subnets_private = {{.AwsCidrSubnetsPrivate}}
aws_cidr_subnets_public = {{.AwsCidrSubnetsPublic}}

aws_bastion_size = "{{.AwsBastionSize}}"
aws_kube_master_num = "{{.AwsKubeMasterNum}}"
aws_kube_master_size = "{{.AwsKubeMasterSize}}"
aws_etcd_num = "{{.AwsEtcdNum}}"

aws_etcd_size = "{{.AwsEtcdSize}}"
aws_kube_worker_num = "{{.AwsKubeWorkerNum}}"
aws_kube_worker_size = "{{.AwsKubeWorkerSize}}"
aws_elb_api_port = "{{.AwsElbAPIPort}}"
k8s_secure_api_port = "{{.K8sSecureAPIPort}}"
kube_insecure_apiserver_address = "{{.KubeInsecureApiserverAddress}}"

default_tags = {
    Env = "devtest"
    Product = "kubernetes"
}
`
