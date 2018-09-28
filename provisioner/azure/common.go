package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Azure struct {
}

func (p Azure) Init(args []string) {
	cluster.NotImplemented()
}

func (p Azure) Setup(args []string) {
	cluster.NotImplemented()
}

func (p Azure) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Azure) Destroy(args []string) {
	cluster.NotImplemented()
}

func NewAzure() cluster.Provisioner {
	provisioner := new(Azure)
	return provisioner
}
