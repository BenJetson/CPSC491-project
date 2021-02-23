FROM golang:1.15

ENV GOPATH /go
ENV GO111MODULE on

RUN curl -sSfL \
    https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin

ENTRYPOINT ["/go/bin/air"]
