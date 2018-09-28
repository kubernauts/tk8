package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type AWS struct {
}

func (p AWS) Init() {
	cluster.KubesprayInit()
	cluster.AWSCreate()
}

func (p AWS) Setup() {
	cluster.AWSInstall()

}

func (p AWS) Upgrade() {
	cluster.NotImplemented()
}

func (p AWS) Destroy() {
	cluster.AWSDestroy()
}

func NewAWS() cluster.Provisioner {
	cluster.SetClusteName()
	provisioner := new(AWS)
	return provisioner
}
