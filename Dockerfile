FROM ubuntu:focal

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update
RUN apt install -yfq make git wget build-essential net-tools
RUN wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz

RUN mkdir -p /root/go/src/github.com/charlesread/gumdrop
ENV GOPATH=/root/go
COPY . /root/go/src/github.com/charlesread/gumdrop

WORKDIR /root/go/src/github.com/charlesread/gumdrop
RUN /usr/local/go/bin/go build gumdrop.go
CMD ./gumdrop
