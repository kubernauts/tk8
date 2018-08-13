# multi-stage build 
## builder stage
FROM golang:alpine as builder
RUN apk --update add make git
COPY ./ /go/src/github.com/kubernauts/tk8
WORKDIR /go/src/github.com/kubernauts/tk8
RUN make bin

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

COPY --from=builder /go/src/github.com/kubernauts/tk8/tk8 /usr/local/bin/tk8

RUN wget https://releases.hashicorp.com/terraform/${TERRVERSION}/terraform_${TERRVERSION}_linux_amd64.zip \
    && unzip terraform_${TERRVERSION}_linux_amd64.zip -d /usr/local/bin/ \
    && rm terraform_${TERRVERSION}_linux_amd64.zip 
#Need git to clone kubespray
#Need '--virtual' packages for proper package import through pip

RUN apk --update add python py-pip openssl ca-certificates py-netaddr git \
    && apk add --virtual build-dependencies \
            python-dev libffi-dev openssl-dev build-base \
    && pip install --upgrade pip cffi ansible \
    && apk add --no-cache openssh \
    && apk del build-dependencies \
    && rm -rf /var/cache/apk/* \
    && chmod +x /usr/local/bin/tk8

# Install requirements for kubectl
RUN apk add -U curl tar gzip

# Install kubectl
RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTLVERSION}/bin/linux/amd64/kubectl && \
  chmod +x /usr/bin/kubectl

CMD [ "/usr/local/bin/tk8" ]
