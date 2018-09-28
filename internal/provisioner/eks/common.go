package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type EKS struct {
}

func (p EKS) Init() {
	cluster.NotImplemented()
}

func (p EKS) Setup() {
	cluster.NotImplemented()
}

func (p EKS) Upgrade() {
	cluster.NotImplemented()
}

func (p EKS) Destroy() {
	cluster.NotImplemented()
}

func NewEKS() cluster.Provisioner {
	provisioner := new(EKS)
	return provisioner
}
