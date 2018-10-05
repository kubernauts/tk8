# Docker-Image

```shell
docker run -it --name tk8 -v ~/.ssh/id_rsa:/root/.ssh/id_rsa kubernautslabs/tk8 [flags][command]
```

## Verwendung des Docker-Image

```shell
docker run -it --name tk8 -v ~/.kube/config:/root/.kube/config kubernautslabs/tk8 [flags][command]
```

## Docker Image pullen

```shell
docker pull kubernautslabs/tk8
```