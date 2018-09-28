package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Azure struct {
}

func (p Azure) Init() {
	cluster.NotImplemented()
}

func (p Azure) Setup() {
	cluster.NotImplemented()
}

func (p Azure) Upgrade() {
	cluster.NotImplemented()
}

func (p Azure) Destroy() {
	cluster.NotImplemented()
}

func NewAzure() cluster.Provisioner {
	provisioner := new(Azure)
	return provisioner
}
