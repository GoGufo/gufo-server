FROM golang:buster AS builder


RUN apt-get update && apt-get install build-essential clang git -y

WORKDIR $GOPATH/src/project/gufo/

COPY . .



ENV CC=clang CGO_ENABLED=1 GOOS=linux GOARCH=amd64


RUN go get -d -v
RUN go build -o /go/bin/gufo gufo.go


FROM ubuntu

COPY --from=builder /go/bin/gufo /go/bin/gufo

WORKDIR /go/bin/

EXPOSE 8090

ENTRYPOINT ["/go/bin/gufo"]
