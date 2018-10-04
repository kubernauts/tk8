# Variables Configuration
variable "cluster-name" {
  default     = "kubernauts"
  type        = "string"
  description = "The name of your EKS Cluster"
}

variable "aws-region" {
  default     = "eu-west-1"
  # availabe regions are:
  # us-east-1 (Virginia)
  # us-west-2 (Oregon)
  # eu-west-1 (Irland)
  type        = "string"
  description = "The AWS Region to deploy EKS"
}

variable "node-instance-type" {
  default     = "t2.medium"
  type        = "string"
  description = "Worker Node EC2 instance type"
}

variable "desired-capacity" {
  default     = 1
  type        = "string"
  description = "Autoscaling Desired node capacity"
}

variable "max-size" {
  default     = 2
  type        = "string"
  description = "Autoscaling maximum node capacity"
}

variable "min-size" {
  default     = 1
  type        = "string"
  description = "Autoscaling Minimum node capacity"
}
variable "public_key_path" {
  default = "~/.ssh/id_rsa.pub"
  type    = "string"
  description = "path to your tk8eks public key"
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