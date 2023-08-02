FROM golang:1.20.1-buster AS builder

RUN apt-get update && apt-get install build-essential clang -y

WORKDIR $GOPATH/src/project/gufo/
COPY . .

RUN go get -u all

ENV CC=clang CGO_ENABLED=1 GOOS=linux GOARCH=amd64

RUN go build -buildmode=plugin -o /var/gufo/lib/login.so plugins/login/*.go
