FROM golang:latest

ENV GOPATH /go

RUN mkdir -p "$GOPATH/src/ec2-info"
ADD . "$GOPATH/src/ec2-info"

WORKDIR $GOPATH/src/ec2-info/app
RUN chmod +x script/*
RUN ./script/build

CMD ./app
