package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Nutanix struct {
}

func (p Nutanix) Init(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Setup(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Destroy(args []string) {
	cluster.NotImplemented()
}

func NewNutanix() cluster.Provisioner {
	provisioner := new(Nutanix)
	return provisioner
}
