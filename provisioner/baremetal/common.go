package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Baremetal struct {
}

func (p Baremetal) Init(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Setup(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Scale(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Reset(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Remove(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Baremetal) Destroy(args []string) {
	cluster.NotImplemented()
}

func NewBaremetal() cluster.Provisioner {
	provisioner := new(Baremetal)
	return provisioner
}
