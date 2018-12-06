# TK8 Provisioner

Tk8 unterstützt verschiedene Platformen zur Bereitstellung eines Kubernetes Clustes.

## Available Provisioner

* [AWS](aws/introduction.md)
* [EKS](aws/introduction.md)
* [Baremetal](baremetal/introduction.md)
* [Openstack](openstack/introduction.md)
* [azure](azure/introduction.md)
* [Nutanix](nutanix/introduction.md)

## Add own Provisioner

Es ist auch möglich einen eigenen Provisioner zu implementieren.
Hier zu erstellen Sie ein neues Repository und fügen eine neue Struktur Ihres Provisioners hinzu. Diese gehört zu dem Package Provisioner. Nun muss die Struktur das interface Provisioner implementieren.

Hier eine Beispiel implementierung

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

Im Packet sollten alle benötigten Scripte enthalten sein um eine Infrastruktur auf der Platform zu erstellen und eine Implementierung für das installieren des Kubernetes Clusters beinhalten.

In manchen Fällen kann die Standard TK8 Implementierung zum installieren von Kubernetes genutzt werden. Hierfür muss eine entsprechende Hosts Datei erzeugt werden und ein Bastion Server in der Infrastruktur existieren.

Erstellen Sie nun ein Issue im Git Repository und verweisen auf die neue Implementierung.

Der Provisioner muss in die Map cmd.provisioners (cmd/provisioners.go) hinzugefügt werden und kann anschließend über die CLI verwendet werden.

```shell
tk8 cluster install MyProvisioner
```

Benötigen Sie Unterstützung oder haben Fragen zur implementierung dann treten Sie unseren [Slack Server](https://kubernauts-slack-join.herokuapp.com/) bei.
