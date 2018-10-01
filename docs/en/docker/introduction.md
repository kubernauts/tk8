#### Using Docker image

```shell
docker run -it --name tk8 -v ~/.ssh/id_rsa:/root/.ssh/id_rsa kubernautslabs/tk8 cluster [flags] [command]


#### Using Docker image

```shell
docker run -it --name tk8 -v ~/.kube/config:/root/.kube/config kubernautslabs/tk8 tk8 addon [flags]
```


### Docker image

```shell
docker pull kubernautslabs/tk8
```