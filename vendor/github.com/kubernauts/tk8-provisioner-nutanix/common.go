package provisioner

import (
	"github.com/kubernauts/tk8-provisioner-nutanix/internal/cluster"
)

type Nutanix struct {
}

var Name string

func (p Nutanix) Init(args []string) {
	Name = cluster.Name
	if len(Name) == 0 {
		Name = "TK8Nutanix"
	}
}

func (p Nutanix) Setup(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Scale(args []string) {
	cluster.NotImplemented()

}

func (p Nutanix) Reset(args []string) {
	cluster.NotImplemented()

}

func (p Nutanix) Remove(args []string) {
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
