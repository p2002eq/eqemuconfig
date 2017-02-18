FROM golang:1.6.2
MAINTAINER Xackery <xackery@gmail.com>

ENV GOPATH /go
ENV USER root

# pre-install known dependencies before the source, so we don't redownload them whenever the source changes

COPY . /go/src/github.com/xackery/eqemuconfig

RUN cd /go/src/github.com/xackery/eqemuconfig \
	&& go get -d -v \
	&& go install \
	&& go test github.com/xackery/eqemuconfig...
