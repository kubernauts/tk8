#
# Workstation External IP

data "http" "workstation-external-ip" {
  url = "http://icanhazip.com"
}

locals {
  workstation-external-cidr = "${chomp(data.http.workstation-external-ip.body)}/32"
}
