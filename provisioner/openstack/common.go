package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Openstack struct {
}

func (p Openstack) Init(args []string) {
	cluster.NotImplemented()
}

func (p Openstack) Setup(args []string) {
	cluster.NotImplemented()
}

func (p Openstack) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Openstack) Destroy(args []string) {
	cluster.NotImplemented()
}

func NewOpenstack() cluster.Provisioner {
	provisioner := new(Openstack)
	return provisioner
}
