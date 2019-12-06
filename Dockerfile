#FROM golang:1.12.8
#GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
FROM alpine:3.10
MAINTAINER oranie

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD . /go/src/myapp
WORKDIR /go/src/myapp
CMD ["/go/src/myapp/main"]