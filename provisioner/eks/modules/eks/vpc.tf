#
# VPC Resources

variable cluster-name {}

variable "aws-region" {}

resource "aws_vpc" "eks" {
  cidr_block = "10.0.0.0/16"

  tags = "${
    map(
     "Name", "${var.cluster-name}-eks-vpc",
     "kubernetes.io/cluster/${var.cluster-name}", "shared",
    )
  }"
}

resource "aws_subnet" "eks" {
  count = 2

  availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = "${aws_vpc.eks.id}"

  tags = "${
    map(
     "Name", "${var.cluster-name}-eks",
     "kubernetes.io/cluster/${var.cluster-name}", "shared",
    )
  }"
}

resource "aws_internet_gateway" "eks" {
  vpc_id = "${aws_vpc.eks.id}"

  tags {
    Name = "${var.cluster-name}-eks-igw"
  }
}

resource "aws_route_table" "eks" {
  vpc_id = "${aws_vpc.eks.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.eks.id}"
  }
}

resource "aws_route_table_association" "eks" {
  count = 2

  subnet_id      = "${aws_subnet.eks.*.id[count.index]}"
  route_table_id = "${aws_route_table.eks.id}"
}
