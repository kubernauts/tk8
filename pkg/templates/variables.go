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

var VariablesCattleEKS = `
variable "service_role" {
  default     = "{{.ServiceRole}}"
  description = "Service linked role for EKS"
  type        = string
}

variable "associate_worker_node_public_ip" {
  default     = {{.AssociatePublicIp}}
  description = "Associate public IP with worker nodes"
  type        = string
}

variable "kubernetes_version" {
  default     = "{{.KubernetesVersion}}"
  description = "Kubernetes version"
  type        = string
}
variable "maximum_nodes" {
  default     = {{.MaximumNodes}}
  description = "Maximum worker nodes"
  type        = string
}

variable "minimum_nodes" {
  default     = {{.MinimumNodes}}
  description = "Minimum size for worker nodes"
  type        = string
}

variable "instance_type" {
  default     = "{{.InstanceType}}"
  description = "Instance type for worker nodes"
  type        = string
}

variable "ami_id" {
  default     = "{{.AmiId}}"
  description = "AMI id for the instance"
  type        = string
}

variable "session_token" {
  default     = "{{.SessionToken}}"
  description = "Session token to use with the client key and secret key if applicable"
  type        = string
}

variable "disk_size" {
  default     = "{{.DiskSize}}"
  description = "Root disk size for instances in GB"
  type        = string
}

variable "rancher_api_url" {
  default     = "{{.RancherApiUrl}}"
  description = "Rancher API URL"
  type        = string
}

variable "rancher_access_key" {
  description = "Rancher server's access key"
}

variable "rancher_secret_key" {
  description = "Rancher server's secret key"
}

variable "rancher_cluster_name" {
  default     = "{{.RancherClusterName}}"
  description = "Rancher cluster name"
  type        = string
}

variable "region" {
  default     = "{{.Region}}"
  description = "AWS region"
  type        = string
}

variable "existing_vpc" {
  default     = {{.ExistingVpc}}
  description = "Use existing VPC for creating clusters"
  type        = string
}

variable "vpc_id" {
  default     = "{{.VpcId}}"
  description = "VPC ID"
  type        = string
}

variable "subnet_id1" {
  default     = "{{.SubnetId1}}"
  description = "subnet id"
  type        = string
}

variable "subnet_id2" {
  default     = "{{.SubnetId2}}"
  description = "Subnet Id 2"
  type        = string
}

variable "subnet_id3" {
  default     = "{{.SubnetId3}}"
  description = "Subnet Id 3"
  type        = string
}

variable "security_group_name" {
  default     = "{{.SecurityGroupName}}"
  description = "security group id"
  type        = string
}

variable "AWS_ACCESS_KEY_ID" {
  description = "AWS access key"
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "AWS secret key"
}

variable "AWS_DEFAULT_REGION" {
  description = "AWS default region"
}
`

var VariablesRKE = `
variable "cluster_name" {
  default     = "{{.ClusterName}}"
  type        = "string"
  description = "The name of your EKS Cluster"
}

variable "authorization" {
  default     = "{{.Authorization}}"
  type        = "string"
  description = "Authorization mode for rke cluster"
}

variable "rke_aws_region" {
  default     = "{{.AWSRegion}}"
  # availabe regions are:
  # us-east-1 (Virginia)
  # us-west-2 (Oregon)
  # eu-west-1 (Irland)
  type        = "string"
  description = "The AWS Region to deploy EKS"
}

variable "rke_node_instance_type" {
  default     = "{{.RKENodeInstanceType}}"
  type        = "string"
  description = "Worker Node EC2 instance type"
}

variable "node_count" {
  default     = {{.NodeCount}}
  type        = "string"
  description = "Autoscaling Desired node capacity"
}
variable "cloud_provider" {
  default = "{{.CloudProvider}}"
  type    = "string"
  description = "cloud provider for rke cluster"
}

variable "AWS_ACCESS_KEY_ID" {
  description = "AWS Access Key"
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "AWS Secret Key"
}

variable "AWS_DEFAULT_REGION" {
  description = "AWS Region"
}
`
var VariablesEKS = `
# Variables Configuration
variable "cluster-name" {
  default     = "{{.ClusterName}}"
  type        = "string"
  description = "The name of your EKS Cluster"
}

variable "aws-region" {
  default     = "{{.AWSRegion}}"
  # availabe regions are:
  # us-east-1 (Virginia)
  # us-west-2 (Oregon)
  # eu-west-1 (Irland)
  type        = "string"
  description = "The AWS Region to deploy EKS"
}

variable "node-instance-type" {
  default     = "{{.NodeInstanceType}}"
  type        = "string"
  description = "Worker Node EC2 instance type"
}

variable "desired-capacity" {
  default     = {{.DesiredCapacity}}
  type        = "string"
  description = "Autoscaling Desired node capacity"
}

variable "max-size" {
  default     = {{.AutoScallingMaxSize}}
  type        = "string"
  description = "Autoscaling maximum node capacity"
}

variable "min-size" {
  default     = {{.AutoScallingMinSize}}
  type        = "string"
  description = "Autoscaling Minimum node capacity"
}
variable "public_key_path" {
  default = "{{.KeyPath}}"
  type    = "string"
  description = "path to your eks public key"
}

variable "AWS_ACCESS_KEY_ID" {
  description = "AWS Access Key"
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "AWS Secret Key"
}

variable "AWS_DEFAULT_REGION" {
  description = "AWS Region"
}
`
var Variables = `
variable "AWS_ACCESS_KEY_ID" {
  description = "AWS Access Key"
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "AWS Secret Key"
}

variable "AWS_SSH_KEY_NAME" {
  description = "Name of the SSH keypair to use in AWS."
}

variable "AWS_DEFAULT_REGION" {
  description = "AWS Region"
}

//General Cluster Settings

variable "aws_cluster_name" {
  description = "Name of AWS Cluster"
}

data "aws_ami" "distro" {
  most_recent = true

  filter {
    name   = "name"
    values = ["{{.OS}}"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["{{.AmiOwner}}"]
}

//AWS VPC Variables

variable "aws_vpc_cidr_block" {
  description = "CIDR Block for VPC"
}

variable "aws_cidr_subnets_private" {
  description = "CIDR Blocks for private subnets in Availability Zones"
  type = "list"
}

variable "aws_cidr_subnets_public" {
  description = "CIDR Blocks for public subnets in Availability Zones"
  type = "list"
}

//AWS EC2 Settings

variable "aws_bastion_size" {
    description = "EC2 Instance Size of Bastion Host"
}

/*
* AWS EC2 Settings
* The number should be divisable by the number of used
* AWS Availability Zones without an remainder.
*/
variable "aws_kube_master_num" {
    description = "Number of Kubernetes Master Nodes"
}

variable "aws_kube_master_size" {
    description = "Instance size of Kube Master Nodes"
}

variable "aws_etcd_num" {
    description = "Number of etcd Nodes"
}

variable "aws_etcd_size" {
    description = "Instance size of etcd Nodes"
}

variable "aws_kube_worker_num" {
    description = "Number of Kubernetes Worker Nodes"
}

variable "aws_kube_worker_size" {
    description = "Instance size of Kubernetes Worker Nodes"
}

/*
* AWS ELB Settings
*
*/
variable "aws_elb_api_port" {
    description = "Port for AWS ELB"
}

variable "k8s_secure_api_port" {
    description = "Secure Port of K8S API Server"
}

variable "default_tags" {
  description = "Default tags for all resources"
  type = "map"
}
`
var VariablesCattleAWS = `
variable "cloudwatch_monitoring" {
  default     = "{{.CloudwatchMonitoring}}"
  description = "Enable/Disable cloudwatch monitoring"
  type        = "string"
}

variable "ami_id" {
  default     = "{{.AmiID}}"
  description = "AMI ID for instances"
  type        = "string"
}

variable "controlplane_instance_type" {
  default     = "{{.ControlPlaneInstanceType}}"
  description = "Control plane instance type"
  type        = "string"
}

variable "request_spot_instances" {
  default     = "{{.RequestSpotInstances}}"
  description = "Request spot instances for the node template"
  type        = "string"
}

variable "spot_price" {
  default     = "{{.SpotPrice}}"
  description = "Spot instances price for the node template"
  type        = "string"
}

variable "root_disk_size" {
  default     = "{{.RootDiskSize}}"
  description = "Root disk size for instances in GB"
  type        = "string"
}

variable "iam_instance_profile_name_worker" {
  default     = "{{.IAMInstanceProfileWorker}}"
  description = "IAM instance profile name for worker"
  type        = "string"
}

variable "iam_instance_profile_name" {
  default     = "{{.IAMInstanceProfile}}"
  description = "IAM instance profile name"
  type        = "string"
}

variable "rancher_api_url" {
  default     = "{{.RancherAPIURL}}"
  description = "Rancher API URL"
  type        = "string"
}

variable "rancher_access_key" {
  description = "Rancher server's access key"
}

variable "rancher_secret_key" {
  description = "Rancher server's secret key"
}

variable "rancher_cluster_name" {
  default     = "{{.RancherClusterName}}"
  description = "Rancher cluster name"
  type        = "string"
}

variable "rke_network_plugin" {
  default     = "{{.RKENetworkPlugin}}"
  description = "Network plugin for cluster"
  type        = "string"
}

variable "region" {
  default     = "{{.Region}}"
  description = "AWS region"
  type        = "string"
}

variable "existing_vpc" {
  default     = {{.ExistingVPC}}
  description = "Use existing VPC for creating clusters"
  type        = "string"
}

variable "vpc_id" {
  default     = "{{.VPCID}}"
  description = "VPC ID"
  type        = "string"
}

variable "subnet_id" {
  default     = "{{.SubnetID}}"
  description = "subnet id"
  type        = "string"
}

variable "security_group_name" {
  default     = "{{.SecurityGroupName}}"
  description = "security group id"
  type        = "string"
}

variable "os" {
  default     = "{{.OS}}"
  description = "ami id - frankfurt"
  type        = "string"
}

variable "worker_instance_type" {
  default     = "{{.WorkerInstanceType}}"
  description = "Instance type"
  type        = "string"
}

variable "overlap_cp_etcd_worker" {
  default     = {{.OverlapCpEtcdWorker}}
  description = "Overlapping planes for node template"
  type        = "string"
}

variable "overlap_node_pool_hostname_prefix" {
  default     = "{{.OverlapHostnamePrefix}}"
  description = "Hostname prefix for overlapped node pools"
  type        = "string"
}

variable "no_overlap_nodepool_master_hostname_prefix" {
  default     = "{{.MasterHostnamePrefix}}"
  description = "Hostname prefix for master node pool"
  type        = "string"
}

variable "no_overlap_nodepool_worker_hostname_prefix" {
  default     = "{{.WorkerHostnamePrefix}}"
  description = "Hostname prefix for worker node pool"
  type        = "string"
}

variable "no_overlap_nodepool_master_quantity" {
  default     = "{{.MasterQuantity}}"
  description = "Node pool master quantity for non-overlapped planes"
  type        = "string"
}

variable "no_overlap_nodepool_worker_quantity" {
  default     = "{{.WorkerQuantity}}"
  description = "Node pool worker quantity for non-overlapped planes"
  type        = "string"
}

variable "overlap_node_pool_quantity" {
  default     = "{{.OverlapQuantity}}"
  description = "Node pool quantity for overlap planes"
  type        = "string"
}

variable "AWS_ACCESS_KEY_ID" {
  description = "AWS access key"
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "AWS secret key"
}

variable "AWS_DEFAULT_REGION" {
  description = "AWS default region"
}
`
