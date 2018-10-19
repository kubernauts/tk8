// Copyright © 2018 The TK8 Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
