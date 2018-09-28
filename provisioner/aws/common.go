package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type AWS struct {
}

func (p AWS) Init(args []string) {
	cluster.KubesprayInit()
	cluster.AWSCreate()
}

func (p AWS) Setup(args []string) {
	cluster.AWSInstall()

}

func (p AWS) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p AWS) Destroy(args []string) {
	cluster.AWSDestroy()
}

func NewAWS() cluster.Provisioner {
	cluster.SetClusteName()
	provisioner := new(AWS)
	return provisioner
}
