# EKS Terraform module

module "eks" {
  source             = "./modules/eks"
  cluster-name       = "${var.cluster-name}"
  aws-region         = "${var.aws-region}"
  node-instance-type = "${var.node-instance-type}"
  desired-capacity   = "${var.desired-capacity}"
  max-size           = "${var.max-size}"
  min-size           = "${var.min-size}"
  public_key_path    = "${var.public_key_path}"
  AWS_ACCESS_KEY_ID    = "${var.AWS_ACCESS_KEY_ID}"
  AWS_SECRET_ACCESS_KEY    = "${var.AWS_SECRET_ACCESS_KEY}"
  AWS_DEFAULT_REGION    = "${var.AWS_DEFAULT_REGION}"
}
