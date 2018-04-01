# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM ubuntu

ARG TERRVERSION=0.11.5

COPY ./tk8 /usr/local/bin/tk8

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367 \
&& apt-get update && apt-get install -y python-pip zip wget git \
&& pip install ansible netaddr \
&& wget https://releases.hashicorp.com/terraform/${TERRVERSION}/terraform_${TERRVERSION}_linux_amd64.zip \
&& unzip terraform_${TERRVERSION}_linux_amd64.zip -d /usr/local/bin/ \
&& rm terraform_${TERRVERSION}_linux_amd64.zip \
&& mkdir /tk8 \
&& chmod +x /usr/local/bin/tk8
