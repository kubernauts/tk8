# TK8 Commissioner

Tk8 supports different platforms to provide a Kubernetes Cluster.

## Avaible Provisioner

* [AWS](aws/introduction.md)
* [EKS](aws/introduction.md)
* [Baremetal](baremetal/introduction.md)
* [Openstack](openstack/introduction.md)
* [azure](azure/introduction.md)
* [Nutanix](nutanix/introduction.md)

## Add own Provisioner

It is also possible to implement your own commissioner.
Here you create a new repository and add a new structure to your commissioner. This belongs to the Commissioner package. Now the structure must implement the interface Commissioner.

Here is an example of implementation

```go
package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type MyProvisioner struct {
}

func (p MyProvisioner) Init(args []string) {
 cluster.NotImplemented()
}

func (p MyProvisioner) Setup(args []string) {
 cluster.NotImplemented()

}

func (p MyProvisioner) Upgrade(args []string) {
 cluster.NotImplemented()
}

func (p MyProvisioner) Destroy(args []string) {
 cluster.NotImplemented()
}

func NewMyProvisioner() cluster.Provisioner {
 cluster.SetClusteName()
 provisioner := new(MyProvisioner)
 return provisioner
}
```

The package should include all scripts needed to create an infrastructure on the platform and an implementation for installing the Kubernetes cluster.

In some cases, the standard TK8 implementation can be used to install Kubernetes. For this, a corresponding Hosts file must be created and a bastion server must exist in the infrastructure.

Now create an issue in the Git Repository and refer to the new implementation.

The provisioner must be added to the map cmd.provisioners (cmd/provisioners.go) and can then be used via the CLI.

```shell
tk8 cluster install MyProvisioner
```

If you need support or have questions about implementation, please join our [Slack Server](https://kubernauts-slack-join.herokuapp.com/).