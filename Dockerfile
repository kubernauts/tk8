# multi-stage build 
## builder stage
FROM golang:alpine as builder
RUN apk --update add make git
COPY ./ /go/src/github.com/kubernauts/tk8
WORKDIR /go/src/github.com/kubernauts/tk8
RUN go get -u . && make bin

## os stage
FROM alpine
#To track exactly which commit is the image built off
ARG VCS_REF=dev
ARG BUILD_DATE=null
#This will be overridden by the build args in hooks folder
ARG TERRVERSION=0.11.7
ARG KUBECTLVERSION=v1.10.5

#Label Schemas to be used for metadata as described at http://label-schema.org/
LABEL  org.label-schema.description="CLI to deploy kubernetes using kubespray and also install additional addons." \
       org.label-schema.usage="docker run kubernauts/tk8:latest [command]" \
       org.label-schema.docker.cmd="docker run kubernauts/tk8:latest [command]" \
       org.label-schema.build-date=$BUILD_DATE \
       org.label-schema.name="kubernauts/tk8" \
       org.label-schema.schema-version="1.0.0-rc.1" \
       org.label-schema.url="https://github.com/kubernauts/tk8" \
       org.label-schema.vcs-ref=$VCS_REF \
       org.label-schema.vcs-url="https://github.com/kubernauts/tk8" \
       org.label-schema.vendor="kubernauts"

COPY --from=builder /go/src/github.com/kubernauts/tk8/tk8ctl /usr/local/bin/tk8ctl

RUN apk --update add \
    python \
    py-pip \
    git \
    gcc \
    ca-certificates \
    py-netaddr \
    python-dev \
    libffi-dev \
    openssl-dev \
    build-base \
    openssh \
    curl \
    tar \
    gzip \
    --virtual build-dependencies \
    --no-cache openssh 

RUN pip install --upgrade pip
RUN pip install --upgrade cffi
RUN pip install --upgrade ansible
RUN pip install --upgrade ansible-modules-hashivault

RUN chmod +x /usr/local/bin/tk8ctl

## Install terraform

RUN wget https://releases.hashicorp.com/terraform/${TERRVERSION}/terraform_${TERRVERSION}_linux_amd64.zip \
    && unzip terraform_${TERRVERSION}_linux_amd64.zip -d /usr/local/bin/ \
    && rm terraform_${TERRVERSION}_linux_amd64.zip 

# Install kubectl
RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTLVERSION}/bin/linux/amd64/kubectl && \
  chmod +x /usr/bin/kubectl

RUN mkdir /tk8

WORKDIR /tk8

ENTRYPOINT [ "/usr/local/bin/tk8ctl" ]
