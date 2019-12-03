#FROM golang:1.12.8
FROM alpine:3.8
MAINTAINER oranie
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD . /go/src/myapp
WORKDIR /go/src/myapp
CMD ["/go/src/myapp/main"]