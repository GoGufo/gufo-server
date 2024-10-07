FROM golang:1.23.0 AS builder


RUN apt-get update && apt-get install build-essential clang git -y

WORKDIR $GOPATH/src/project/gufo/

COPY . .



ENV CC=clang CGO_ENABLED=1 GOOS=linux GOARCH=amd64



RUN go build -o /go/bin/gufo gufo.go


FROM ubuntu

ADD var/ /var/gufo/
COPY --from=builder /go/bin/gufo /go/bin/gufo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/project/gufo/config/settings.toml /var/gufo/config/

WORKDIR /go/bin/

EXPOSE 8090

ENTRYPOINT ["/go/bin/gufo"]
