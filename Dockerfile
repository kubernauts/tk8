FROM golang:alpine as builder

COPY ./ /go/src/tk8

WORKDIR /go/src/tk8

RUN apk add --update git curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    && dep ensure -v

RUN go install tk8

FROM alpine

ARG TERRVERSION=0.11.7

COPY --from=builder /go/bin/tk8 /usr/local/bin/tk8

RUN wget https://releases.hashicorp.com/terraform/${TERRVERSION}/terraform_${TERRVERSION}_linux_amd64.zip \
    && unzip terraform_${TERRVERSION}_linux_amd64.zip -d /usr/local/bin/ \
    && rm terraform_${TERRVERSION}_linux_amd64.zip 

RUN apk add --no-cache py-netaddr ansible \
    && chmod +x /usr/local/bin/tk8

CMD [ "--help" ]

ENTRYPOINT [ "/usr/local/bin/tk8" ]
