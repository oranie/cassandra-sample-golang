#GOOS=linux GOARCH=amd64 CGO_ENABLED=0
#go build main.go

#FROM amazonlinux
FROM alpine:3.10
MAINTAINER oranie

#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD . /go/src/myapp
WORKDIR /go/src/myapp
RUN chmod 755 ./main ./run.sh
ENTRYPOINT ["/go/src/myapp/main"]
CMD ["/go/src/myapp/run.sh && echo 'success'"]
