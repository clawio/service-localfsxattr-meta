FROM golang:1.5
MAINTAINER Hugo González Labrador

ENV CLAWIO_LOCALFSXATTR_META_DATADIR /tmp/localfsxattr-mysql
ENV CLAWIO_LOCALFSXATTR_META_TMPDIR /tmp/localfsxattr-mysql
ENV CLAWIO_LOCALFSXATTR_META_PORT 57011
ENV CLAWIO_LOCALFSXATTR_META_LOGLEVEL "error"
ENV CLAWIO_LOCALFSXATTR_META_PROP "service-localfsxattr-mysqlprop:57013"
ENV CLAWIO_LOCALFSXATTR_META_PROPMAXACTIVE 1024
ENV CLAWIO_LOCALFSXATTR_META_PROPMAXIDLE 1024
ENV CLAWIO_LOCALFSXATTR_META_PROPMAXCONCURRENCY 1024
ENV CLAWIO_SHAREDSECRET secret

ADD . /go/src/github.com/clawio/service-localfsxattr-meta
WORKDIR /go/src/github.com/clawio/service-localfsxattr-meta

RUN go get -u github.com/tools/godep
RUN godep restore
RUN go install

ENTRYPOINT /go/bin/service-localfsxattr-meta
EXPOSE 57011

