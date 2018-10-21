# Managing Cluster Lifecycle

There will be times when you would want to scale up or down a cluster, upgrade it to a new version, or just completely reset it and remove Kubernetes from it. We already does support some of these functionality in the TK8 cli, whereas rest are already under development. You can easily scale up the cluster, or reset it removing Kubernetes from the infrastructure.


## Scaling the cluster

To scale the cluster, you can increase the desired count of any of the following in the [config.yaml](../../../../config.yaml.example) to meet your requirements:

* `aws_kube_master_num`
* `aws_kube_worker_num`
* `aws_etcd_num`

Make sure you are in the same directory where you executed `tk8 cluster install aws` with the inventory directory.
If you use a different workspace name with the --name flag please provide it on scaling too.

To scale the provisioned cluster run:

```shell
tk8 cluster scale aws
```

Once executed a confirmation would be needed to overwrite the existing inventory file post which above command will scale your infrastructure as well as the Kubernetes cluster together to the specified capacity.

## Reset the cluster

Make sure you are in the same directory where you executed `tk8 cluster install aws` with the inventory directory.
If you use a different workspace name with the --name flag please provide it on resetting too.

To reset the provisioned cluster run:

```shell
tk8 cluster reset aws
```

Once executed the current kubernetes installation get removed and a new setup will run.

## Remove the cluster

Make sure you are in the same directory where you executed `tk8 cluster install aws` with the inventory directory.
If you use a different workspace name with the --name flag please provide it on resetting too.

To reset the provisioned cluster run:

```shell
tk8 cluster remove aws
```

Once executed the current kubernetes installation get removed from the infrastructure.
