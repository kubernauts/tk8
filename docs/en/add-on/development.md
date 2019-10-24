# TK8 Add-On Development

The Add-On implementation needs to be a  general solution so we provide a command to deliver a default example add-on which can be customized by the user.

## Development commands

The command to create the example add-on gets created with this.

```shell
tk8ctl addon create my-addon
```

This command pulls the tk8-addon-develop from GitHub and creates a new folder below ./addons/my-addon

The example is a simple nginx deployment and a LoadBalancer service to expose this. So the user who creates this add-on can directly use it and apply it to the k8s cluster

```shell
tk8ctl addon install my-addon
```

and could remove it from the k8s cluster with

```shell
tk8ctl addon destroy my-addon
```

The default developer add-on doesn't contain a main.sh file. But we need to create documentation for it. Our own add-ons could use and need it.

## TK8 Add-on structure

For a general use of Add-Ons with git we defined a standard frame which contains the folder structure, the yml structure and an example.

The Folder structure

→ addons

| → my-addon

| →  | → LICENCE

| →  | → Readme.md

| →  | → main.yml

| →  | → main.sh

The main.yml contains all the needed information for k8s and will create all the deployments and services which are needed.

Optionally there is a main.sh which can used to download external repositories or to create a main.yml