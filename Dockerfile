FROM golang:1.5
MAINTAINER Hugo González Labrador

ENV CLAWIO_LOCALSTOREXATTRMETA_DATADIR /tmp
ENV CLAWIO_LOCALSTOREXATTRMETA_TMPDIR /tmp
ENV CLAWIO_LOCALSTOREXATTRMETA_PORT 57011
ENV CLAWIO_LOCALSTOREXATTRMETA_PROP "service-localstorexattr-prop:57013"
ENV CLAWIO_SHAREDSECRET secret

ADD . /go/src/github.com/clawio/service.localstorexattr.meta
WORKDIR /go/src/github.com/clawio/service.localstorexattr.meta

RUN go get -u github.com/tools/godep
RUN godep restore
RUN go install

ENTRYPOINT /go/bin/service.localstorexattr.meta

EXPOSE 57011

