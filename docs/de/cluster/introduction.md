# Cluster Dokumentation

## CLI Help

```shell
Used to initialize both the kubespray repo and the specified public infrastructure

Usage:
  tk8 cluster [flags]
  tk8 cluster [command]

Available Commands:
  destroy     destroy the infrastructure
  install     install the infrastructure
  upgrade     Manages the infrastructure on AWS

Flags:
  -h, --help   help for cluster

Use "tk8 cluster [command] --help" for more information about a command.
```

## Befehle

Der Befehl Cluster verfügt über verschiedene Unterbefehle um ein Cluster auf einer gewünschten Platform zu erstellen und Kubernetes darauf zu installieren.
Ebenso liefert er Möglichlkeiten um ein Cluster zu entfernen oder zu scallieren.

### cluster install

```shell
Usage:
  tk8 cluster install [aws|azure|baremetal|eks|nutanix|openstack] [flags]

Flags:
  -h, --help          help for install
      --name string   name of the cluster workspace (default "kubernauts")
```

```shell
tk8 cluster install *provisioner*
tk8 cluster install *provisioner* -n my-cluster -r eu-west-1 -k myKeyPairName -o ubuntu
```

### Cluster destroy

```shell
Usage:
  tk8 cluster destroy [aws|azure|baremetal|eks|nutanix|openstack] [flags]

Flags:
  -h, --help          help for destroy
      --name string   name of the cluster workspace (default "kubernauts")
```

```shell
tk8 cluster destroy *provisioner*
tk8 cluster destroy *provisioner* -n my-cluster
```

### Cluster update

```shell
tk8 cluster update *provisioner*
tk8 cluster update *provisioner* -n my-cluster -v 1.11.2
```

### Cluster scale

```shell
tk8 cluster scale *provisioner*
```

### Cluster reset

```shell
tk8 cluster reset *provisioner*
```

### Cluster remove

```shell
tk8 cluster remove *provisioner*
```