package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Baremetal struct {
}

func (p Baremetal) Init() {
	cluster.NotImplemented()

}

func (p Baremetal) Setup() {
	cluster.NotImplemented()

}

func (p Baremetal) Upgrade() {
	cluster.NotImplemented()
}

func (p Baremetal) Destroy() {
	cluster.NotImplemented()
}

func NewBaremetal() cluster.Provisioner {
	provisioner := new(Baremetal)
	return provisioner
}
