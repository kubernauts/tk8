package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Openstack struct {
}

func (p Openstack) Init() {
	cluster.NotImplemented()
}

func (p Openstack) Setup() {
	cluster.NotImplemented()
}

func (p Openstack) Upgrade() {
	cluster.NotImplemented()
}

func (p Openstack) Destroy() {
	cluster.NotImplemented()
}

func NewOpenstack() cluster.Provisioner {
	provisioner := new(Openstack)
	return provisioner
}
