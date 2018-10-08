# TK8 Add-On

## Using

How to install an add-on with tk8?

We need to store every Add-On in a single repository on Github for public add-ons and on gitlab for internal add-ons.

We make a switch with install tk8 add-ons and provide a shortcut. A shortcut could also be a local add-on so we need to check first if there is one on in the folder. If not, check if there a tk8-addon- on GitHub.

### Use the complete Path

```shell
tk8 addon install https://github.com/kubernauts/tk-addon-rancher
tk8 addon install https://github.com/kubernauts/tk-addon-prometheus
tk8 addon install https://github.com/kubernauts/tk-addon-grafana
tk8 addon install https://github.com/kubernauts/tk-addon-monitoring-stack
tk8 addon install https://github.com/kubernauts/tk-addon-elk
tk8 addon install https://github.com/kubernauts/tk-addon-...
tk8 addon install https://github.com/USERNAME/ADDON-REPO
```

### Use the shortcut

```shell
tk8 addon install rancher
tk8 addon install prometheus
tk8 addon install grafana
tk8 addon install monitoring-stack
tk8 addon install elk
```

### Destroy a add-on

```shell
tk8 addon destroy rancher
tk8 addon destroy prometheus
tk8 addon destroy grafana
tk8 addon destroy monitoring-stack
tk8 addon destroy elk
```

## Development

Create an add-on.
The create method of tk8 creates a new add-on in the local folder. This add-on is a simple example and provides all that we need to work with this add-on.

[More information here](development.md)

```shell
tk8 addon create my-addon
```