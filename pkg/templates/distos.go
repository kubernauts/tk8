package templates

var DistVariablesRKE = `
data "aws_ami" "distro" {
  most_recent = true

  filter {
    name   = "name"
    values = ["{{.NodeOS}}"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["{{.AmiOwner}}"]
}

variable "ssh_user" {
  default = "{{.User}}"
}
`
