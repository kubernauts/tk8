package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Nutanix struct {
}

func (p Nutanix) Init() {
	cluster.NotImplemented()
}

func (p Nutanix) Setup() {
	cluster.NotImplemented()
}

func (p Nutanix) Upgrade() {
	cluster.NotImplemented()
}

func (p Nutanix) Destroy() {
	cluster.NotImplemented()
}

func NewNutanix() cluster.Provisioner {
	provisioner := new(Nutanix)
	return provisioner
}
