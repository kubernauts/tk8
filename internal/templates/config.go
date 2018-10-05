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

// Config is used to generate the default config.
var Config = `
aws:
   clustername: {{.ClusterName}}
   os: centos
   aws_access_key_id: {{.AccessKey}}
   aws_secret_access_key: {{.SecretKey}}
   aws_ssh_keypair: {{.SSHName}}
   aws_default_region: us-east-1
   aws_vpc_cidr_block : "10.250.192.0/18"
   aws_cidr_subnets_private : '["10.250.192.0/20","10.250.208.0/20"]'
   aws_cidr_subnets_public : '["10.250.224.0/20","10.250.240.0/20"]'
   aws_bastion_size : "t2.medium"
   aws_kube_master_num : 1
   aws_kube_master_size : "t2.medium"
   aws_etcd_num : 1
   aws_etcd_size : "t2.medium"
   aws_kube_worker_num : 1
   aws_kube_worker_size : "t2.medium"
   aws_elb_api_port : 6443
   k8s_secure_api_port : 6443
   kube_insecure_apiserver_address : "0.0.0.0"
   kubeadm_enabled: false
   `
